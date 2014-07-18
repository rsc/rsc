package main

// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
/*
Stack layout parameters.
Included both by runtime (compiled via 6c) and linkers (compiled via gcc).

The per-goroutine g->stackguard is set to point StackGuard bytes
above the bottom of the stack.  Each function compares its stack
pointer against g->stackguard to check for overflow.  To cut one
instruction from the check sequence for functions with tiny frames,
the stack is allowed to protrude StackSmall bytes below the stack
guard.  Functions with large frames don't bother with the check and
always call morestack.  The sequences are (for amd64, others are
similar):

	guard = g->stackguard
	frame = function's stack frame size
	argsize = size of function arguments (call + return)

	stack frame size <= StackSmall:
		CMPQ guard, SP
		JHI 3(PC)
		MOVQ m->morearg, $(argsize << 32)
		CALL morestack(SB)

	stack frame size > StackSmall but < StackBig
		LEAQ (frame-StackSmall)(SP), R0
		CMPQ guard, R0
		JHI 3(PC)
		MOVQ m->morearg, $(argsize << 32)
		CALL morestack(SB)

	stack frame size >= StackBig:
		MOVQ m->morearg, $((argsize << 32) | frame)
		CALL morestack(SB)

The bottom StackGuard - StackSmall bytes are important: there has
to be enough room to execute functions that refuse to check for
stack overflow, either because they need to be adjacent to the
actual caller's frame (deferproc) or because they handle the imminent
stack overflow (morestack).

For example, deferproc might call malloc, which does one of the
above checks (without allocating a full frame), which might trigger
a call to morestack.  This sequence needs to fit in the bottom
section of the stack.  On amd64, morestack's frame is 40 bytes, and
deferproc's frame is 56 bytes.  That fits well within the
StackGuard - StackSmall = 128 bytes at the bottom.
The linkers explore all possible call traces involving non-splitting
functions to make sure that this limit cannot be violated.
*/
const (
	StackSystem_stack        = 0
	StackExtra_stack         = 2048
	StackMin_stack           = 8192
	StackSystemRounded_stack = StackSystem_stack + (-StackSystem_stack & (StackMin_stack - 1))
	FixedStack_stack         = StackMin_stack + StackSystemRounded_stack
	StackBig_stack           = 4096
	StackGuard_stack         = 256 + StackSystem_stack
	StackSmall_stack         = 128
	StackLimit_stack         = StackGuard_stack - StackSystem_stack - StackSmall_stack
	StackTop_stack           = 88
)

// Goroutine preemption request.
// Stored into g->stackguard0 to cause split stack check failure.
// Must be greater than any real sp.
// 0xfffffade in hex.
const (
	StackPreempt_stack = -1314
)
