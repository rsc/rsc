package main

const (
	StackSystem_stack        = 0
	StackExtra_stack         = 2048
	StackMin_stack           = 8192
	StackSystemRounded_stack = StackSystem_stack + (-StackSystem_stack & (StackMin_stack - 1))
	FixedStack               = StackMin_stack + StackSystemRounded_stack
	StackBig_stack           = 4096
	StackGuard_stack         = 256 + StackSystem_stack
	StackSmall_stack         = 128
	StackLimit_stack         = StackGuard_stack - StackSystem_stack - StackSmall_stack
	StackTop_stack           = 88
	StackPreempt_stack       = -1314
)

const (
	NOPROF_textflag   = 1
	DUPOK_textflag    = 2
	NOSPLIT_textflag  = 4
	RODATA_textflag   = 8
	NOPTR_textflag    = 16
	WRAPPER_textflag  = 32
	NEEDCTXT_textflag = 64
)
