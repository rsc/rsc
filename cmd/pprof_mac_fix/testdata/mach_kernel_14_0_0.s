_version:
	.ascii "Darwin Kernel Version 14.0.0: Sat Sep 27 03:58:47 PDT 2014; root:xnu-2782.1.97~11/RELEASE_X86_64\0"

.globl _main
_main:
	ret

.globl _psignal_internal
_psignal_internal:
	ret

.globl _task_vtimer_clear
_task_vtimer_clear:
	ret

.globl _task_vtimer_set
_task_vtimer_set:
	ret
.globl OSCompareAndSwap
OSCompareAndSwap:
	ret

.globl addupc_task
addupc_task:
	ret

.globl bsdinit_task
bsdinit_task:
	ret

.globl current_proc
current_proc:
	ret

.globl get_useraddr
get_useraddr:
	ret

.globl hw_lock_unlock
hw_lock_unlock:
	ret

.globl issignal_locked
issignal_locked:
	ret

.globl itimerdecr
itimerdecr:
	ret

.globl lck_mtx_lock
lck_mtx_lock:
	ret

.globl lck_mtx_unlock
lck_mtx_unlock:
	ret

.globl postsig_locked
postsig_locked:
	ret

.globl proc_findinternal
proc_findinternal:
	ret

.globl proc_rele
proc_rele:
	ret

.globl task_resume_internal
task_resume_internal:
	ret

.globl task_suspend_internal
task_suspend_internal:
	ret

.globl task_vtimer_set
task_vtimer_set:
	ret

.globl usimple_lock
usimple_lock:
	ret

.globl _current_thread
_current_thread:
//   0xffffff8000418e90 <+0>:	push   %rbp
	.byte 0x55;
//   0xffffff8000418e91 <+1>:	mov    %rsp,%rbp
	.byte 0x48; .byte 0x89; .byte 0xe5;
//   0xffffff8000418e94 <+4>:	mov    %gs:0x8,%rax
	.byte 0x65; .byte 0x48; .byte 0x8b; .byte 0x04; .byte 0x25; .byte 0x08; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff8000418e9d <+13>:	pop    %rbp
	.byte 0x5d;
//   0xffffff8000418e9e <+14>:	retq   
	.byte 0xc3;
//   0xffffff8000418e9f <+15>:	nop
	.byte 0x90;

.globl _bsd_ast
.align 16
_bsd_ast:
//   0xffffff80007dbe20 <+0>:	push   %rbp
	.byte 0x55;
//   0xffffff80007dbe21 <+1>:	mov    %rsp,%rbp
	.byte 0x48; .byte 0x89; .byte 0xe5;
//   0xffffff80007dbe24 <+4>:	push   %r15
	.byte 0x41; .byte 0x57;
//   0xffffff80007dbe26 <+6>:	push   %r14
	.byte 0x41; .byte 0x56;
//   0xffffff80007dbe28 <+8>:	push   %r13
	.byte 0x41; .byte 0x55;
//   0xffffff80007dbe2a <+10>:	push   %r12
	.byte 0x41; .byte 0x54;
//   0xffffff80007dbe2c <+12>:	push   %rbx
	.byte 0x53;
//   0xffffff80007dbe2d <+13>:	push   %rax
	.byte 0x50;
//   0xffffff80007dbe2e <+14>:	mov    %rdi,%rbx
	.byte 0x48; .byte 0x89; .byte 0xfb;
//   0xffffff80007dbe31 <+17>:	callq  0xffffff8000859a30 <current_proc>
	call current_proc
//   0xffffff80007dbe36 <+22>:	mov    %rax,%r15
	.byte 0x49; .byte 0x89; .byte 0xc7;
//   0xffffff80007dbe39 <+25>:	test   %r15,%r15
	.byte 0x4d; .byte 0x85; .byte 0xff;
//   0xffffff80007dbe3c <+28>:	je     0xffffff80007dc2b1 <bsd_ast+1169>
	.byte 0x0f; .byte 0x84; .byte 0x6f; .byte 0x04; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe42 <+34>:	mov    0x3c8(%rbx),%r12
	.byte 0x4c; .byte 0x8b; .byte 0xa3; .byte 0xc8; .byte 0x03; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe49 <+41>:	mov    $0x8020,%eax
	.byte 0xb8; .byte 0x20; .byte 0x80; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe4e <+46>:	and    0x160(%r15),%eax
	.byte 0x41; .byte 0x23; .byte 0x87; .byte 0x60; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe55 <+53>:	cmp    $0x8020,%eax
	.byte 0x3d; .byte 0x20; .byte 0x80; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe5a <+58>:	jne    0xffffff80007dbe96 <bsd_ast+118>
	.byte 0x75; .byte 0x3a;
//   0xffffff80007dbe5c <+60>:	lea    0x160(%r15),%rbx
	.byte 0x49; .byte 0x8d; .byte 0x9f; .byte 0x60; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe63 <+67>:	callq  0xffffff800041b9b0 <get_useraddr>
	call get_useraddr
//   0xffffff80007dbe68 <+72>:	mov    %eax,%esi
	.byte 0x89; .byte 0xc6;
//   0xffffff80007dbe6a <+74>:	mov    $0x1,%edx
	.byte 0xba; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe6f <+79>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dbe72 <+82>:	callq  0xffffff80007edb50 <addupc_task>
	call addupc_task
//   0xffffff80007dbe77 <+87>:	nopw   0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x0f; .byte 0x1f; .byte 0x84; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbe80 <+96>:	mov    (%rbx),%edi
	.byte 0x8b; .byte 0x3b;
//   0xffffff80007dbe82 <+98>:	mov    %edi,%esi
	.byte 0x89; .byte 0xfe;
//   0xffffff80007dbe84 <+100>:	and    $0xffff7fff,%esi
	.byte 0x81; .byte 0xe6; .byte 0xff; .byte 0x7f; .byte 0xff; .byte 0xff;
//   0xffffff80007dbe8a <+106>:	mov    %rbx,%rdx
	.byte 0x48; .byte 0x89; .byte 0xda;
//   0xffffff80007dbe8d <+109>:	callq  0xffffff800089a4f0 <OSCompareAndSwap>
	call OSCompareAndSwap
//   0xffffff80007dbe92 <+114>:	test   %al,%al
	.byte 0x84; .byte 0xc0;
//   0xffffff80007dbe94 <+116>:	je     0xffffff80007dbe80 <bsd_ast+96>
	.byte 0x74; .byte 0xea;
//   0xffffff80007dbe96 <+118>:	mov    $0xffffffff,%r13d
	.byte 0x41; .byte 0xbd; .byte 0xff; .byte 0xff; .byte 0xff; .byte 0xff;
//   0xffffff80007dbe9c <+124>:	cmpq   $0x0,0x1c8(%r15)
	.byte 0x49; .byte 0x83; .byte 0xbf; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbea4 <+132>:	jne    0xffffff80007dbeb4 <bsd_ast+148>
	.byte 0x75; .byte 0x0e;
//   0xffffff80007dbea6 <+134>:	cmpl   $0x0,0x1d0(%r15)
	.byte 0x41; .byte 0x83; .byte 0xbf; .byte 0xd0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbeae <+142>:	je     0xffffff80007dbf70 <bsd_ast+336>
	.byte 0x0f; .byte 0x84; .byte 0xbc; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbeb4 <+148>:	lea    0x1b8(%r15),%rsi
	.byte 0x49; .byte 0x8d; .byte 0xb7; .byte 0xb8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbebb <+155>:	mov    %gs:0x8,%rdx
	.byte 0x65; .byte 0x48; .byte 0x8b; .byte 0x14; .byte 0x25; .byte 0x08; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbec4 <+164>:	cmpl   $0x0,0x1b8(%rdx)
	.byte 0x83; .byte 0xba; .byte 0xb8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbecb <+171>:	mov    0x1e0(%rdx),%rax
	.byte 0x48; .byte 0x8b; .byte 0x82; .byte 0xe0; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbed2 <+178>:	je     0xffffff80007dbedd <bsd_ast+189>
	.byte 0x74; .byte 0x09;
//   0xffffff80007dbed4 <+180>:	mov    0x1c8(%rdx),%rcx
	.byte 0x48; .byte 0x8b; .byte 0x8a; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbedb <+187>:	jmp    0xffffff80007dbee4 <bsd_ast+196>
	.byte 0xeb; .byte 0x07;
//   0xffffff80007dbedd <+189>:	mov    0x158(%rdx),%rcx
	.byte 0x48; .byte 0x8b; .byte 0x8a; .byte 0x58; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbee4 <+196>:	mov    %rcx,0x1e0(%rdx)
	.byte 0x48; .byte 0x89; .byte 0x8a; .byte 0xe0; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbeeb <+203>:	sub    %eax,%ecx
	.byte 0x29; .byte 0xc1;
//   0xffffff80007dbeed <+205>:	and    %r13,%rcx
	.byte 0x4c; .byte 0x21; .byte 0xe9;
//   0xffffff80007dbef0 <+208>:	mov    %rcx,%rax
	.byte 0x48; .byte 0x89; .byte 0xc8;
//   0xffffff80007dbef3 <+211>:	shr    $0x9,%rax
	.byte 0x48; .byte 0xc1; .byte 0xe8; .byte 0x09;
//   0xffffff80007dbef7 <+215>:	movabs $0x44b82fa09b5a53,%rdx
	.byte 0x48; .byte 0xba; .byte 0x53; .byte 0x5a; .byte 0x9b; .byte 0xa0; .byte 0x2f; .byte 0xb8; .byte 0x44; .byte 0x00;
//   0xffffff80007dbf01 <+225>:	mul    %rdx
	.byte 0x48; .byte 0xf7; .byte 0xe2;
//   0xffffff80007dbf04 <+228>:	shr    $0xb,%rdx
	.byte 0x48; .byte 0xc1; .byte 0xea; .byte 0x0b;
//   0xffffff80007dbf08 <+232>:	imul   $0x3b9aca00,%rdx,%rax
	.byte 0x48; .byte 0x69; .byte 0xc2; .byte 0x00; .byte 0xca; .byte 0x9a; .byte 0x3b;
//   0xffffff80007dbf0f <+239>:	sub    %rax,%rcx
	.byte 0x48; .byte 0x29; .byte 0xc1;
//   0xffffff80007dbf12 <+242>:	shr    $0x3,%rcx
	.byte 0x48; .byte 0xc1; .byte 0xe9; .byte 0x03;
//   0xffffff80007dbf16 <+246>:	movabs $0x20c49ba5e353f7cf,%rdx
	.byte 0x48; .byte 0xba; .byte 0xcf; .byte 0xf7; .byte 0x53; .byte 0xe3; .byte 0xa5; .byte 0x9b; .byte 0xc4; .byte 0x20;
//   0xffffff80007dbf20 <+256>:	mov    %rcx,%rax
	.byte 0x48; .byte 0x89; .byte 0xc8;
//   0xffffff80007dbf23 <+259>:	mul    %rdx
	.byte 0x48; .byte 0xf7; .byte 0xe2;
//   0xffffff80007dbf26 <+262>:	shr    $0x4,%rdx
	.byte 0x48; .byte 0xc1; .byte 0xea; .byte 0x04;
//   0xffffff80007dbf2a <+266>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dbf2d <+269>:	callq  0xffffff80007e9310 <itimerdecr>
	call itimerdecr
//   0xffffff80007dbf32 <+274>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
//   0xffffff80007dbf34 <+276>:	jne    0xffffff80007dbf70 <bsd_ast+336>
	.byte 0x75; .byte 0x3a;
//   0xffffff80007dbf36 <+278>:	cmpq   $0x0,0x1c8(%r15)
	.byte 0x49; .byte 0x83; .byte 0xbf; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf3e <+286>:	jne    0xffffff80007dbf4e <bsd_ast+302>
	.byte 0x75; .byte 0x0e;
//   0xffffff80007dbf40 <+288>:	cmpl   $0x0,0x1d0(%r15)
	.byte 0x41; .byte 0x83; .byte 0xbf; .byte 0xd0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf48 <+296>:	je     0xffffff80007dc2c0 <bsd_ast+1184>
	.byte 0x0f; .byte 0x84; .byte 0x72; .byte 0x03; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf4e <+302>:	mov    0x18(%r15),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7f; .byte 0x18;
//   0xffffff80007dbf52 <+306>:	mov    $0x1,%esi
	.byte 0xbe; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf57 <+311>:	callq  0xffffff800035e620 <task_vtimer_set>
	call task_vtimer_set
//   0xffffff80007dbf5c <+316>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
//   0xffffff80007dbf5e <+318>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
//   0xffffff80007dbf60 <+320>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
//   0xffffff80007dbf62 <+322>:	mov    $0x1a,%r8d
	.byte 0x41; .byte 0xb8; .byte 0x1a; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf68 <+328>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dbf6b <+331>:	callq  0xffffff80007dc590 <_psignal_internal>
	call _psignal_internal
//   0xffffff80007dbf70 <+336>:	cmpq   $0x0,0x1e8(%r15)
	.byte 0x49; .byte 0x83; .byte 0xbf; .byte 0xe8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf78 <+344>:	jne    0xffffff80007dbf88 <bsd_ast+360>
	.byte 0x75; .byte 0x0e;
//   0xffffff80007dbf7a <+346>:	cmpl   $0x0,0x1f0(%r15)
	.byte 0x41; .byte 0x83; .byte 0xbf; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf82 <+354>:	je     0xffffff80007dc03e <bsd_ast+542>
	.byte 0x0f; .byte 0x84; .byte 0xb6; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf88 <+360>:	lea    0x1d8(%r15),%rsi
	.byte 0x49; .byte 0x8d; .byte 0xb7; .byte 0xd8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf8f <+367>:	mov    %gs:0x8,%rdi
	.byte 0x65; .byte 0x48; .byte 0x8b; .byte 0x3c; .byte 0x25; .byte 0x08; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf98 <+376>:	mov    0x158(%rdi),%rbx
	.byte 0x48; .byte 0x8b; .byte 0x9f; .byte 0x58; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbf9f <+383>:	add    0x1c8(%rdi),%rbx
	.byte 0x48; .byte 0x03; .byte 0x9f; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbfa6 <+390>:	mov    0x1e8(%rdi),%eax
	.byte 0x8b; .byte 0x87; .byte 0xe8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbfac <+396>:	mov    %ebx,%ecx
	.byte 0x89; .byte 0xd9;
//   0xffffff80007dbfae <+398>:	sub    %eax,%ecx
	.byte 0x29; .byte 0xc1;
//   0xffffff80007dbfb0 <+400>:	and    %r13,%rcx
	.byte 0x4c; .byte 0x21; .byte 0xe9;
//   0xffffff80007dbfb3 <+403>:	mov    %rcx,%rax
	.byte 0x48; .byte 0x89; .byte 0xc8;
//   0xffffff80007dbfb6 <+406>:	shr    $0x9,%rax
	.byte 0x48; .byte 0xc1; .byte 0xe8; .byte 0x09;
//   0xffffff80007dbfba <+410>:	movabs $0x44b82fa09b5a53,%rdx
	.byte 0x48; .byte 0xba; .byte 0x53; .byte 0x5a; .byte 0x9b; .byte 0xa0; .byte 0x2f; .byte 0xb8; .byte 0x44; .byte 0x00;
//   0xffffff80007dbfc4 <+420>:	mul    %rdx
	.byte 0x48; .byte 0xf7; .byte 0xe2;
//   0xffffff80007dbfc7 <+423>:	shr    $0xb,%rdx
	.byte 0x48; .byte 0xc1; .byte 0xea; .byte 0x0b;
//   0xffffff80007dbfcb <+427>:	imul   $0x3b9aca00,%rdx,%rax
	.byte 0x48; .byte 0x69; .byte 0xc2; .byte 0x00; .byte 0xca; .byte 0x9a; .byte 0x3b;
//   0xffffff80007dbfd2 <+434>:	sub    %rax,%rcx
	.byte 0x48; .byte 0x29; .byte 0xc1;
//   0xffffff80007dbfd5 <+437>:	shr    $0x3,%rcx
	.byte 0x48; .byte 0xc1; .byte 0xe9; .byte 0x03;
//   0xffffff80007dbfd9 <+441>:	movabs $0x20c49ba5e353f7cf,%rdx
	.byte 0x48; .byte 0xba; .byte 0xcf; .byte 0xf7; .byte 0x53; .byte 0xe3; .byte 0xa5; .byte 0x9b; .byte 0xc4; .byte 0x20;
//   0xffffff80007dbfe3 <+451>:	mov    %rcx,%rax
	.byte 0x48; .byte 0x89; .byte 0xc8;
//   0xffffff80007dbfe6 <+454>:	mul    %rdx
	.byte 0x48; .byte 0xf7; .byte 0xe2;
//   0xffffff80007dbfe9 <+457>:	shr    $0x4,%rdx
	.byte 0x48; .byte 0xc1; .byte 0xea; .byte 0x04;
//   0xffffff80007dbfed <+461>:	test   %edx,%edx
	.byte 0x85; .byte 0xd2;
//   0xffffff80007dbfef <+463>:	je     0xffffff80007dbff8 <bsd_ast+472>
	.byte 0x74; .byte 0x07;
//   0xffffff80007dbff1 <+465>:	mov    %rbx,0x1e8(%rdi)
	.byte 0x48; .byte 0x89; .byte 0x9f; .byte 0xe8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dbff8 <+472>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dbffb <+475>:	callq  0xffffff80007e9310 <itimerdecr>
	call itimerdecr
//   0xffffff80007dc000 <+480>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
//   0xffffff80007dc002 <+482>:	jne    0xffffff80007dc03e <bsd_ast+542>
	.byte 0x75; .byte 0x3a;
//   0xffffff80007dc004 <+484>:	cmpq   $0x0,0x1e8(%r15)
	.byte 0x49; .byte 0x83; .byte 0xbf; .byte 0xe8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc00c <+492>:	jne    0xffffff80007dc01c <bsd_ast+508>
	.byte 0x75; .byte 0x0e;
//   0xffffff80007dc00e <+494>:	cmpl   $0x0,0x1f0(%r15)
	.byte 0x41; .byte 0x83; .byte 0xbf; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc016 <+502>:	je     0xffffff80007dc2e1 <bsd_ast+1217>
	.byte 0x0f; .byte 0x84; .byte 0xc5; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc01c <+508>:	mov    0x18(%r15),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7f; .byte 0x18;
//   0xffffff80007dc020 <+512>:	mov    $0x2,%esi
	.byte 0xbe; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc025 <+517>:	callq  0xffffff800035e620 <task_vtimer_set>
	call task_vtimer_set
//   0xffffff80007dc02a <+522>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
//   0xffffff80007dc02c <+524>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
//   0xffffff80007dc02e <+526>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
//   0xffffff80007dc030 <+528>:	mov    $0x1b,%r8d
	.byte 0x41; .byte 0xb8; .byte 0x1b; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc036 <+534>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dc039 <+537>:	callq  0xffffff80007dc590 <_psignal_internal>
	call _psignal_internal
//   0xffffff80007dc03e <+542>:	cmpq   $0x0,0x1f8(%r15)
	.byte 0x49; .byte 0x83; .byte 0xbf; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc046 <+550>:	jne    0xffffff80007dc056 <bsd_ast+566>
	.byte 0x75; .byte 0x0e;
//   0xffffff80007dc048 <+552>:	cmpl   $0x0,0x200(%r15)
	.byte 0x41; .byte 0x83; .byte 0xbf; .byte 0x00; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc050 <+560>:	je     0xffffff80007dc15c <bsd_ast+828>
	.byte 0x0f; .byte 0x84; .byte 0x06; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc056 <+566>:	mov    %gs:0x8,%rax
	.byte 0x65; .byte 0x48; .byte 0x8b; .byte 0x04; .byte 0x25; .byte 0x08; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc05f <+575>:	mov    0x158(%rax),%rdx
	.byte 0x48; .byte 0x8b; .byte 0x90; .byte 0x58; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc066 <+582>:	add    0x1c8(%rax),%rdx
	.byte 0x48; .byte 0x03; .byte 0x90; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc06d <+589>:	mov    0x1f0(%rax),%esi
	.byte 0x8b; .byte 0xb0; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc073 <+595>:	mov    %edx,%ecx
	.byte 0x89; .byte 0xd1;
//   0xffffff80007dc075 <+597>:	sub    %esi,%ecx
	.byte 0x29; .byte 0xf1;
//   0xffffff80007dc077 <+599>:	mov    %rdx,0x1f0(%rax)
	.byte 0x48; .byte 0x89; .byte 0x90; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc07e <+606>:	and    %r13,%rcx
	.byte 0x4c; .byte 0x21; .byte 0xe9;
//   0xffffff80007dc081 <+609>:	mov    %rcx,%rax
	.byte 0x48; .byte 0x89; .byte 0xc8;
//   0xffffff80007dc084 <+612>:	shr    $0x9,%rax
	.byte 0x48; .byte 0xc1; .byte 0xe8; .byte 0x09;
//   0xffffff80007dc088 <+616>:	movabs $0x44b82fa09b5a53,%rdx
	.byte 0x48; .byte 0xba; .byte 0x53; .byte 0x5a; .byte 0x9b; .byte 0xa0; .byte 0x2f; .byte 0xb8; .byte 0x44; .byte 0x00;
//   0xffffff80007dc092 <+626>:	mul    %rdx
	.byte 0x48; .byte 0xf7; .byte 0xe2;
//   0xffffff80007dc095 <+629>:	shr    $0xb,%rdx
	.byte 0x48; .byte 0xc1; .byte 0xea; .byte 0x0b;
//   0xffffff80007dc099 <+633>:	imul   $0x3b9aca00,%rdx,%rax
	.byte 0x48; .byte 0x69; .byte 0xc2; .byte 0x00; .byte 0xca; .byte 0x9a; .byte 0x3b;
//   0xffffff80007dc0a0 <+640>:	sub    %rax,%rcx
	.byte 0x48; .byte 0x29; .byte 0xc1;
//   0xffffff80007dc0a3 <+643>:	shr    $0x3,%rcx
	.byte 0x48; .byte 0xc1; .byte 0xe9; .byte 0x03;
//   0xffffff80007dc0a7 <+647>:	movabs $0x20c49ba5e353f7cf,%rdx
	.byte 0x48; .byte 0xba; .byte 0xcf; .byte 0xf7; .byte 0x53; .byte 0xe3; .byte 0xa5; .byte 0x9b; .byte 0xc4; .byte 0x20;
//   0xffffff80007dc0b1 <+657>:	mov    %rcx,%rax
	.byte 0x48; .byte 0x89; .byte 0xc8;
//   0xffffff80007dc0b4 <+660>:	mul    %rdx
	.byte 0x48; .byte 0xf7; .byte 0xe2;
//   0xffffff80007dc0b7 <+663>:	mov    %rdx,%rbx
	.byte 0x48; .byte 0x89; .byte 0xd3;
//   0xffffff80007dc0ba <+666>:	shr    $0x4,%rbx
	.byte 0x48; .byte 0xc1; .byte 0xeb; .byte 0x04;
//   0xffffff80007dc0be <+670>:	lea    0x108(%r15),%r14
	.byte 0x4d; .byte 0x8d; .byte 0xb7; .byte 0x08; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc0c5 <+677>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc0c8 <+680>:	callq  0xffffff8000416160 <usimple_lock>
	call usimple_lock
//   0xffffff80007dc0cd <+685>:	mov    0x1f8(%r15),%rcx
	.byte 0x49; .byte 0x8b; .byte 0x8f; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc0d4 <+692>:	test   %rcx,%rcx
	.byte 0x48; .byte 0x85; .byte 0xc9;
//   0xffffff80007dc0d7 <+695>:	mov    0x200(%r15),%eax
	.byte 0x41; .byte 0x8b; .byte 0x87; .byte 0x00; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc0de <+702>:	jg     0xffffff80007dc133 <bsd_ast+787>
	.byte 0x7f; .byte 0x53;
//   0xffffff80007dc0e0 <+704>:	cmp    %ebx,%eax
	.byte 0x39; .byte 0xd8;
//   0xffffff80007dc0e2 <+706>:	jg     0xffffff80007dc133 <bsd_ast+787>
	.byte 0x7f; .byte 0x4f;
//   0xffffff80007dc0e4 <+708>:	movl   $0x0,0x200(%r15)
	.byte 0x41; .byte 0xc7; .byte 0x87; .byte 0x00; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc0ef <+719>:	movq   $0x0,0x1f8(%r15)
	.byte 0x49; .byte 0xc7; .byte 0x87; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc0fa <+730>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc0fd <+733>:	callq  0xffffff8000411ea0 <hw_lock_unlock>
	call hw_lock_unlock
//   0xffffff80007dc102 <+738>:	mov    0x18(%r15),%rbx
	.byte 0x49; .byte 0x8b; .byte 0x5f; .byte 0x18;
//   0xffffff80007dc106 <+742>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc109 <+745>:	callq  0xffffff8000412500 <lck_mtx_lock>
	call lck_mtx_lock
//   0xffffff80007dc10e <+750>:	andb   $0xfb,0xc0(%rbx)
	.byte 0x80; .byte 0xa3; .byte 0xc0; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0xfb;
//   0xffffff80007dc115 <+757>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc118 <+760>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc11d <+765>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
//   0xffffff80007dc11f <+767>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
//   0xffffff80007dc121 <+769>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
//   0xffffff80007dc123 <+771>:	mov    $0x18,%r8d
	.byte 0x41; .byte 0xb8; .byte 0x18; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc129 <+777>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dc12c <+780>:	callq  0xffffff80007dc590 <_psignal_internal>
	call _psignal_internal
//   0xffffff80007dc131 <+785>:	jmp    0xffffff80007dc15c <bsd_ast+828>
	.byte 0xeb; .byte 0x29;
//   0xffffff80007dc133 <+787>:	sub    %ebx,%eax
	.byte 0x29; .byte 0xd8;
//   0xffffff80007dc135 <+789>:	mov    %eax,0x200(%r15)
	.byte 0x41; .byte 0x89; .byte 0x87; .byte 0x00; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc13c <+796>:	jns    0xffffff80007dc154 <bsd_ast+820>
	.byte 0x79; .byte 0x16;
//   0xffffff80007dc13e <+798>:	dec    %rcx
	.byte 0x48; .byte 0xff; .byte 0xc9;
//   0xffffff80007dc141 <+801>:	mov    %rcx,0x1f8(%r15)
	.byte 0x49; .byte 0x89; .byte 0x8f; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc148 <+808>:	add    $0xf4240,%eax
	.byte 0x05; .byte 0x40; .byte 0x42; .byte 0x0f; .byte 0x00;
//   0xffffff80007dc14d <+813>:	mov    %eax,0x200(%r15)
	.byte 0x41; .byte 0x89; .byte 0x87; .byte 0x00; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc154 <+820>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc157 <+823>:	callq  0xffffff8000411ea0 <hw_lock_unlock>
	call hw_lock_unlock
//   0xffffff80007dc15c <+828>:	movzbl 0x239(%r12),%r8d
	.byte 0x45; .byte 0x0f; .byte 0xb6; .byte 0x84; .byte 0x24; .byte 0x39; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc165 <+837>:	test   %r8d,%r8d
	.byte 0x45; .byte 0x85; .byte 0xc0;
//   0xffffff80007dc168 <+840>:	je     0xffffff80007dc181 <bsd_ast+865>
	.byte 0x74; .byte 0x17;
//   0xffffff80007dc16a <+842>:	movb   $0x0,0x239(%r12)
	.byte 0x41; .byte 0xc6; .byte 0x84; .byte 0x24; .byte 0x39; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc173 <+851>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
//   0xffffff80007dc175 <+853>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
//   0xffffff80007dc177 <+855>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
//   0xffffff80007dc179 <+857>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dc17c <+860>:	callq  0xffffff80007dc590 <_psignal_internal>
	call _psignal_internal
//   0xffffff80007dc181 <+865>:	cmpb   $0x0,0x238(%r12)
	.byte 0x41; .byte 0x80; .byte 0xbc; .byte 0x24; .byte 0x38; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc18a <+874>:	je     0xffffff80007dc1ba <bsd_ast+922>
	.byte 0x74; .byte 0x2e;
//   0xffffff80007dc18c <+876>:	movb   $0x0,0x238(%r12)
	.byte 0x41; .byte 0xc6; .byte 0x84; .byte 0x24; .byte 0x38; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc195 <+885>:	lea    0x58(%r15),%rbx
	.byte 0x49; .byte 0x8d; .byte 0x5f; .byte 0x58;
//   0xffffff80007dc199 <+889>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc19c <+892>:	callq  0xffffff8000412500 <lck_mtx_lock>
	call lck_mtx_lock
//   0xffffff80007dc1a1 <+897>:	movb   $0x1,0x278(%r15)
	.byte 0x41; .byte 0xc6; .byte 0x87; .byte 0x78; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x01;
//   0xffffff80007dc1a9 <+905>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc1ac <+908>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc1b1 <+913>:	mov    0x18(%r15),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7f; .byte 0x18;
//   0xffffff80007dc1b5 <+917>:	callq  0xffffff800035d250 <task_suspend_internal>
	call task_suspend_internal
//   0xffffff80007dc1ba <+922>:	mov    0x230(%r12),%rdi
	.byte 0x49; .byte 0x8b; .byte 0xbc; .byte 0x24; .byte 0x30; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc1c2 <+930>:	test   %rdi,%rdi
	.byte 0x48; .byte 0x85; .byte 0xff;
//   0xffffff80007dc1c5 <+933>:	je     0xffffff80007dc223 <bsd_ast+1027>
	.byte 0x74; .byte 0x5c;
//   0xffffff80007dc1c7 <+935>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
//   0xffffff80007dc1c9 <+937>:	callq  0xffffff80007cd810 <proc_findinternal>
	call proc_findinternal
//   0xffffff80007dc1ce <+942>:	mov    %rax,%r14
	.byte 0x49; .byte 0x89; .byte 0xc6;
//   0xffffff80007dc1d1 <+945>:	movq   $0x0,0x230(%r12)
	.byte 0x49; .byte 0xc7; .byte 0x84; .byte 0x24; .byte 0x30; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc1dd <+957>:	test   %r14,%r14
	.byte 0x4d; .byte 0x85; .byte 0xf6;
//   0xffffff80007dc1e0 <+960>:	je     0xffffff80007dc223 <bsd_ast+1027>
	.byte 0x74; .byte 0x41;
//   0xffffff80007dc1e2 <+962>:	lea    0x58(%r14),%rbx
	.byte 0x49; .byte 0x8d; .byte 0x5e; .byte 0x58;
//   0xffffff80007dc1e6 <+966>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc1e9 <+969>:	callq  0xffffff8000412500 <lck_mtx_lock>
	call lck_mtx_lock
//   0xffffff80007dc1ee <+974>:	cmpb   $0x0,0x278(%r14)
	.byte 0x41; .byte 0x80; .byte 0xbe; .byte 0x78; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc1f6 <+982>:	je     0xffffff80007dc213 <bsd_ast+1011>
	.byte 0x74; .byte 0x1b;
//   0xffffff80007dc1f8 <+984>:	movb   $0x0,0x278(%r14)
	.byte 0x41; .byte 0xc6; .byte 0x86; .byte 0x78; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc200 <+992>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc203 <+995>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc208 <+1000>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
//   0xffffff80007dc20c <+1004>:	callq  0xffffff800035d310 <task_resume_internal>
	call task_resume_internal
//   0xffffff80007dc211 <+1009>:	jmp    0xffffff80007dc21b <bsd_ast+1019>
	.byte 0xeb; .byte 0x08;
//   0xffffff80007dc213 <+1011>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
//   0xffffff80007dc216 <+1014>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc21b <+1019>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc21e <+1022>:	callq  0xffffff80007cd790 <proc_rele>
	call proc_rele
//   0xffffff80007dc223 <+1027>:	lea    0x58(%r15),%r14
	.byte 0x4d; .byte 0x8d; .byte 0x77; .byte 0x58;
//   0xffffff80007dc227 <+1031>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc22a <+1034>:	callq  0xffffff8000412500 <lck_mtx_lock>
	call lck_mtx_lock
//   0xffffff80007dc22f <+1039>:	mov    %gs:0x8,%rax
	.byte 0x65; .byte 0x48; .byte 0x8b; .byte 0x04; .byte 0x25; .byte 0x08; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc238 <+1048>:	testb  $0x1,0x394(%rax)
	.byte 0xf6; .byte 0x80; .byte 0x94; .byte 0x03; .byte 0x00; .byte 0x00; .byte 0x01;
//   0xffffff80007dc23f <+1055>:	je     0xffffff80007dc293 <bsd_ast+1139>
	.byte 0x74; .byte 0x52;
//   0xffffff80007dc241 <+1057>:	mov    0x11c(%r12),%eax
	.byte 0x41; .byte 0x8b; .byte 0x84; .byte 0x24; .byte 0x1c; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc249 <+1065>:	mov    0x124(%r12),%ecx
	.byte 0x41; .byte 0x8b; .byte 0x8c; .byte 0x24; .byte 0x24; .byte 0x01; .byte 0x00; .byte 0x00;
//   0xffffff80007dc251 <+1073>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
//   0xffffff80007dc253 <+1075>:	testb  $0x4,0x165(%r15)
	.byte 0x41; .byte 0xf6; .byte 0x87; .byte 0x65; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x04;
//   0xffffff80007dc25b <+1083>:	jne    0xffffff80007dc264 <bsd_ast+1092>
	.byte 0x75; .byte 0x07;
//   0xffffff80007dc25d <+1085>:	mov    0x2c4(%r15),%edx
	.byte 0x41; .byte 0x8b; .byte 0x97; .byte 0xc4; .byte 0x02; .byte 0x00; .byte 0x00;
//   0xffffff80007dc264 <+1092>:	or     %edx,%ecx
	.byte 0x09; .byte 0xd1;
//   0xffffff80007dc266 <+1094>:	or     $0x10100,%ecx
	.byte 0x81; .byte 0xc9; .byte 0x00; .byte 0x01; .byte 0x01; .byte 0x00;
//   0xffffff80007dc26c <+1100>:	xor    $0xfffefeff,%ecx
	.byte 0x81; .byte 0xf1; .byte 0xff; .byte 0xfe; .byte 0xfe; .byte 0xff;
//   0xffffff80007dc272 <+1106>:	test   %ecx,%eax
	.byte 0x85; .byte 0xc8;
//   0xffffff80007dc274 <+1108>:	jne    0xffffff80007dc287 <bsd_ast+1127>
	.byte 0x75; .byte 0x11;
//   0xffffff80007dc276 <+1110>:	jmp    0xffffff80007dc293 <bsd_ast+1139>
	.byte 0xeb; .byte 0x1b;
//   0xffffff80007dc278 <+1112>:	nopl   0x0(%rax,%rax,1)
	.byte 0x0f; .byte 0x1f; .byte 0x84; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
//   0xffffff80007dc280 <+1120>:	mov    %eax,%edi
	.byte 0x89; .byte 0xc7;
//   0xffffff80007dc282 <+1122>:	callq  0xffffff80007dba00 <postsig_locked>
	call postsig_locked
//   0xffffff80007dc287 <+1127>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
//   0xffffff80007dc28a <+1130>:	callq  0xffffff80007dadd0 <issignal_locked>
	call issignal_locked
//   0xffffff80007dc28f <+1135>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
//   0xffffff80007dc291 <+1137>:	jne    0xffffff80007dc280 <bsd_ast+1120>
	.byte 0x75; .byte 0xed;
//   0xffffff80007dc293 <+1139>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc296 <+1142>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc29b <+1147>:	mov    other_sym(%rip),%al        # 0xffffff8000ac8f90
	mov    other_sym(%rip),%al 
//   0xffffff80007dc2a1 <+1153>:	test   %al,%al
	.byte 0x84; .byte 0xc0;
//   0xffffff80007dc2a3 <+1155>:	jne    0xffffff80007dc2b1 <bsd_ast+1169>
	.byte 0x75; .byte 0x0c;
//   0xffffff80007dc2a5 <+1157>:	movb   $0x1,other_sym(%rip)        # 0xffffff8000ac8f90
	movb   $0x1,other_sym(%rip) 
//   0xffffff80007dc2ac <+1164>:	callq  0xffffff80007992c0 <bsdinit_task>
	call bsdinit_task
//   0xffffff80007dc2b1 <+1169>:	add    $0x8,%rsp
	.byte 0x48; .byte 0x83; .byte 0xc4; .byte 0x08;
//   0xffffff80007dc2b5 <+1173>:	pop    %rbx
	.byte 0x5b;
//   0xffffff80007dc2b6 <+1174>:	pop    %r12
	.byte 0x41; .byte 0x5c;
//   0xffffff80007dc2b8 <+1176>:	pop    %r13
	.byte 0x41; .byte 0x5d;
//   0xffffff80007dc2ba <+1178>:	pop    %r14
	.byte 0x41; .byte 0x5e;
//   0xffffff80007dc2bc <+1180>:	pop    %r15
	.byte 0x41; .byte 0x5f;
//   0xffffff80007dc2be <+1182>:	pop    %rbp
	.byte 0x5d;
//   0xffffff80007dc2bf <+1183>:	retq   
	.byte 0xc3;
//   0xffffff80007dc2c0 <+1184>:	mov    0x18(%r15),%r14
	.byte 0x4d; .byte 0x8b; .byte 0x77; .byte 0x18;
//   0xffffff80007dc2c4 <+1188>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc2c7 <+1191>:	callq  0xffffff8000412500 <lck_mtx_lock>
	call lck_mtx_lock
//   0xffffff80007dc2cc <+1196>:	andb   $0xfe,0xc0(%r14)
	.byte 0x41; .byte 0x80; .byte 0xa6; .byte 0xc0; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0xfe;
//   0xffffff80007dc2d4 <+1204>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc2d7 <+1207>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc2dc <+1212>:	jmpq   0xffffff80007dbf5c <bsd_ast+316>
	.byte 0xe9; .byte 0x7b; .byte 0xfc; .byte 0xff; .byte 0xff;
//   0xffffff80007dc2e1 <+1217>:	mov    0x18(%r15),%r14
	.byte 0x4d; .byte 0x8b; .byte 0x77; .byte 0x18;
//   0xffffff80007dc2e5 <+1221>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc2e8 <+1224>:	callq  0xffffff8000412500 <lck_mtx_lock>
	call lck_mtx_lock
//   0xffffff80007dc2ed <+1229>:	andb   $0xfd,0xc0(%r14)
	.byte 0x41; .byte 0x80; .byte 0xa6; .byte 0xc0; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0xfd;
//   0xffffff80007dc2f5 <+1237>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
//   0xffffff80007dc2f8 <+1240>:	callq  0xffffff8000412b00 <lck_mtx_unlock>
	call lck_mtx_unlock
//   0xffffff80007dc2fd <+1245>:	jmpq   0xffffff80007dc02a <bsd_ast+522>
	.byte 0xe9; .byte 0x28; .byte 0xfd; .byte 0xff; .byte 0xff;
//   0xffffff80007dc302 <+1250>:	data32 data32 data32 data32 nopw %cs:0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x66; .byte 0x66; .byte 0x66; .byte 0x66; .byte 0x2e; .byte 0x0f; .byte 0x1f; .byte 0x84; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;

.globl after_bsd_ast
after_bsd_ast:
	ret

.data
.globl other_sym
other_sym:
	.byte 0x80
