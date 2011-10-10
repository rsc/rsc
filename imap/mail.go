package imap

import (
	"fmt"
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
	FlagInbox  // Gmail extension
	FlagAllMail  // Gmail extension
	FlagDrafts  // Gmail extension
	FlagSent  // Gmail extension
	FlagSpam  // Gmail extension
	FlagStarred  // Gmail extension
	FlagTrash  // Gmail extension
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
	Name string  // name of mailbox
	Elem string  // last element in name
	
	client *Client
	parent *Box  // parent in hierarchy
	child []*Box  // child boxes
	dead bool  // box no longer exists
	inbox bool  // box is inbox
	flags Flags  // allowed flags
	permFlags Flags  // client-modifiable permanent flags
	readOnly bool  // box is read-only
	exists int64  // number of messages in box (according to server)
	maxSeen int64  // maximum message number seen (why?)
	unseen int64  // number of first unseen message
	validity int64  // UID validity base number
	msgByUID map[int64]*Msg
	nextUID int64  // the next UID we expect to see (for polling)
}

func (c *Client) newBox(name, sep string, inbox bool) *Box {
//	c.dlk.Lock()
//	defer c.dlk.Unlock()
	if b := c.boxByName[name]; b != nil {
		return b
	}
	
	b := &Box{
		Name: name,
		Elem: name,
		inbox: inbox,
		client: c,
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
	box *Box  // box containing message
	number int64  // message number in box (up to date?)
	uid int64  // unique id of message on server
	date int64  // date (seconds)
	flags Flags  // message flags
	size int64  // size in bytes
	lines int64  // number of lines
	hdr *Hdr  // MIME header

	root MsgPart  // top-level message part

	gmailMsgid int64  // gmail message id
	gmailThrid int64  // gmail thread id
}

func (b *Box) newMsg(uid int64) *Msg {
//	b.client.dlk.Lock()
//	defer b.client.dlk.Unlock()
	if m := b.msgByUID[uid]; m != nil {
		return m
	}
	if b.msgByUID == nil {
		b.msgByUID = map[int64]*Msg{}
	}
	m := &Msg{
		uid: uid,
		box: b,
	}
	m.root.msg = m
	b.msgByUID[uid] = m
	return m
}

// TODO: Hdr->MsgHdr

// A Hdr represents a message header.
type Hdr struct {
	Date string
	Subject string
	From []Addr
	Sender []Addr
	ReplyTo []Addr
	To []Addr
	CC []Addr
	BCC []Addr
	InReplyTo string
	MessageID string
	Digest string
}

// An Addr represents a single, named email address.
// If Name is empty, only the email address is known.
// If Email is empty, the Addr represents an unspecified (but named) group.
type Addr struct {
	Name string
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
	msg *Msg
	id string
	mimeType string
	contentID string
	desc string
	encoding string
	size int64
	lines int64
	charset string
	name string
	hdr *Hdr
	
	raw []byte  // raw message
	rawHeader []byte  // raw RFC-2822 header, for message/rfc822
	rawBody []byte  // raw RFC-2822 body, for message/rfc822
	mimeHeader []byte  // mime header, for attachments

	child []*MsgPart
}

func (p *MsgPart) newPart() *MsgPart {
	dot := "."
	if p.id == "" {  // no dot at root
		dot = ""
	}
	pp := &MsgPart{
		msg: p.msg,
		id: fmt.Sprint(p.id, dot, 1+len(p.child)),
	}
	p.child = append(p.child, pp)
	return pp
}

func (p *MsgPart) text() []byte {
	return decodeText(p.raw, p.encoding, p.charset)
}
