all: \
	mach_kernel_10_8_0 \
	mach_kernel_10_8_0_fix \
	mach_kernel_11_4_2 \
	mach_kernel_11_4_2_fix \
	mach_kernel_12_4_0 \
	mach_kernel_12_4_0_fix \

%: %.s
	gcc -o $* $*.s

