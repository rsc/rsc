package imap

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Flags uint32

const (
	FlagJunk Flags = 1 << iota
	FlagNonJunk
	FlagReplied
	FlagFlagged
	FlagDeleted
	FlagDraft
	FlagRecent
	FlagSeen
	FlagNoInferiors
	FlagNoSelect
	FlagMarked
	FlagUnMarked
	FlagHasChildren
	FlagHasNoChildren
	FlagInbox     // Gmail extension
	FlagAllMail   // Gmail extension
	FlagDrafts    // Gmail extension
	FlagSent      // Gmail extension
	FlagSpam      // Gmail extension
	FlagStarred   // Gmail extension
	FlagTrash     // Gmail extension
	FlagImportant // Gmail extension
)

var flagNames = []string{
	"Junk",
	"NonJunk",
	"\\Answered",
	"\\Flagged",
	"\\Deleted",
	"\\Draft",
	"\\Recent",
	"\\Seen",
	"\\NoInferiors",
	"\\NoSelect",
	"\\Marked",
	"\\UnMarked",
	"\\HasChildren",
	"\\HasNoChildren",
	"\\Inbox",
	"\\AllMail",
	"\\Drafts",
	"\\Sent",
	"\\Spam",
	"\\Starred",
	"\\Trash",
	"\\Important",
}

// A Box represents an IMAP mailbox.
type Box struct {
	Name   string // name of mailbox
	Elem   string // last element in name
	Client *Client

	parent    *Box   // parent in hierarchy
	child     []*Box // child boxes
	dead      bool   // box no longer exists
	inbox     bool   // box is inbox
	flags     Flags  // allowed flags
	permFlags Flags  // client-modifiable permanent flags
	readOnly  bool   // box is read-only
	exists    int64  // number of messages in box (according to server)
	maxSeen   int64  // maximum message number seen (why?)
	unseen    int64  // number of first unseen message
	validity  int64  // UID validity base number
	msgByUID  map[uint64]*Msg
	nextUID   int64 // the next UID we expect to see (for polling)
}

func (c *Client) Boxes() []*Box {
	c.data.lock()
	defer c.data.unlock()

	box := make([]*Box, len(c.allBox))
	copy(box, c.allBox)
	return box
}

func (c *Client) Box(name string) *Box {
	c.data.lock()
	defer c.data.unlock()

	return c.boxByName[name]
}

func (c *Client) Inbox() *Box {
	c.data.lock()
	defer c.data.unlock()

	return c.inbox
}

func (c *Client) newBox(name, sep string, inbox bool) *Box {
	c.data.mustBeLocked()
	if b := c.boxByName[name]; b != nil {
		return b
	}

	b := &Box{
		Name:    name,
		Elem:    name,
		Client:  c,
		inbox:   inbox,
		nextUID: 1,
	}
	if !inbox {
		b.parent = c.rootBox
	}
	if !inbox && sep != "" && name != c.root {
		if i := strings.LastIndex(name, sep); i >= 0 {
			b.Elem = name[i+len(sep):]
			b.parent = c.newBox(name[:i], sep, false)
		}
	}
	c.allBox = append(c.allBox, b)
	c.boxByName[name] = b
	if b.parent != nil {
		b.parent.child = append(b.parent.child, b)
	}
	return b
}

// A Msg represents an IMAP message.
type Msg struct {
	Box         *Box    // box containing message
	Date        int64   // date (seconds)
	Flags       Flags   // message flags
	Bytes       int64   // size in bytes
	Lines       int64   // number of lines
	Hdr         *MsgHdr // MIME header
	Root        MsgPart // top-level message part
	GmailID     uint64  // Gmail message id
	GmailThread uint64  // Gmail thread id
	UID         uint64  // unique id for this message

	deleted bool
	number  int64 // message number in box (up to date?)
}

// TODO: Return os.Error too

type byUID []*Msg

func (x byUID) Len() int           { return len(x) }
func (x byUID) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byUID) Less(i, j int) bool { return x[i].UID < x[j].UID }

func (b *Box) Msgs() []*Msg {
	b.Client.data.lock()
	defer b.Client.data.unlock()

	msgs := make([]*Msg, len(b.msgByUID))
	n := 0
	for _, m := range b.msgByUID {
		msgs[n] = m
		n++
	}
	sort.Sort(byUID(msgs))
	return msgs
}

func (b *Box) newMsg(uid uint64) *Msg {
	b.Client.data.mustBeLocked()
	if m := b.msgByUID[uid]; m != nil {
		return m
	}
	if b.msgByUID == nil {
		b.msgByUID = map[uint64]*Msg{}
	}
	m := &Msg{
		UID: uid,
		Box: b,
	}
	m.Root.Msg = m
	b.msgByUID[uid] = m
	return m
}

func (b *Box) Delete(msgs []*Msg) os.Error {
	for _, m := range msgs {
		if m.Box != b {
			return fmt.Errorf("messages not from this box")
		}
	}
	b.Client.io.lock()
	defer b.Client.io.unlock()
	err := b.Client.deleteList(msgs)
	if err == nil {
		b.Client.data.lock()
		defer b.Client.data.unlock()
		for _, m := range msgs {
			if m.Flags&FlagDeleted != 0 {
				b.msgByUID[m.UID] = nil, false
			}
		}
	}
	return err
}

func (m *Msg) Deleted() bool {
	// Racy but okay.  Can add a lock later if it matters.
	return m.Flags&FlagDeleted != 0
}

// A Hdr represents a message header.
type MsgHdr struct {
	Date      string
	Subject   string
	From      []Addr
	Sender    []Addr
	ReplyTo   []Addr
	To        []Addr
	CC        []Addr
	BCC       []Addr
	InReplyTo string
	MessageID string
	Digest    string
}

// An Addr represents a single, named email address.
// If Name is empty, only the email address is known.
// If Email is empty, the Addr represents an unspecified (but named) group.
type Addr struct {
	Name  string
	Email string
}

func (a Addr) String() string {
	if a.Email == "" {
		return a.Name
	}
	if a.Name == "" {
		return a.Email
	}
	return a.Name + " <" + a.Email + ">"
}

// A MsgPart represents a single part of a MIME-encoded message.
type MsgPart struct {
	Msg       *Msg // containing message
	Type      string
	ContentID string
	Desc      string
	Encoding  string
	Bytes     int64
	Lines     int64
	Charset   string
	Name      string
	Hdr       *MsgHdr
	ID        string
	Child     []*MsgPart

	raw        []byte // raw message
	rawHeader  []byte // raw RFC-2822 header, for message/rfc822
	rawBody    []byte // raw RFC-2822 body, for message/rfc822
	mimeHeader []byte // mime header, for attachments
}

func (p *MsgPart) newPart() *MsgPart {
	p.Msg.Box.Client.data.mustBeLocked()
	dot := "."
	if p.ID == "" { // no dot at root
		dot = ""
	}
	pp := &MsgPart{
		Msg: p.Msg,
		ID:  fmt.Sprint(p.ID, dot, 1+len(p.Child)),
	}
	p.Child = append(p.Child, pp)
	return pp
}

func (p *MsgPart) Text() []byte {
	c := p.Msg.Box.Client
	var raw []byte
	c.data.lock()
	if p == &p.Msg.Root {
		raw = p.rawBody
		c.data.unlock()
		if raw == nil {
			c.io.lock()
			if raw = p.rawBody; raw == nil {
				c.fetch(p, "TEXT")
				raw = p.rawBody
			}
			c.io.unlock()
		}
	} else {
		raw = p.raw
		c.data.unlock()
		if raw == nil {
			c.io.lock()
			if raw = p.raw; raw == nil {
				c.fetch(p, "")
				raw = p.raw
			}
			c.io.unlock()
		}
	}
	return decodeText(raw, p.Encoding, p.Charset, false)
}

func (p *MsgPart) Raw() []byte {
	c := p.Msg.Box.Client
	var raw []byte
	c.data.lock()
	raw = p.rawBody
	c.data.unlock()
	if raw == nil {
		c.io.lock()
		if raw = p.rawBody; raw == nil {
			c.fetch(p, "")
			raw = p.rawBody
		}
		c.io.unlock()
	}
	return raw
}
