package ptx

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"sync"
	"unsafe"
)

// pointers passed to CGO must be kept alive manually
// so we keep then here.
// TODO: how about one struct inside the func. will leak not so much and be parallelizeable.
var (
	reducedot_lock        sync.Mutex
	reducedot_code        cu.Function
	reducedot_stream      cu.Stream
	reducedot_arg_x1      cu.DevicePtr
	reducedot_arg_x2      cu.DevicePtr
	reducedot_arg_dst     cu.DevicePtr
	reducedot_arg_initVal float32
	reducedot_arg_n       int

	reducedot_argptr = [...]unsafe.Pointer{
		unsafe.Pointer(&reducedot_arg_x1),
		unsafe.Pointer(&reducedot_arg_x2),
		unsafe.Pointer(&reducedot_arg_dst),
		unsafe.Pointer(&reducedot_arg_initVal),
		unsafe.Pointer(&reducedot_arg_n)}
)

// CUDA kernel wrapper for reducedot.
// The kernel is launched in a separate stream so that it can be parallel with memcpys etc.
// The stream is synchronized before this call returns.
func K_reducedot(x1 cu.DevicePtr, x2 cu.DevicePtr, dst cu.DevicePtr, initVal float32, n int, gridDim, blockDim cu.Dim3) {
	reducedot_lock.Lock()

	if reducedot_stream == 0 {
		reducedot_stream = cu.StreamCreate()
		//core.Log("Loading PTX code for reducedot")
		reducedot_code = cu.ModuleLoadData(reducedot_ptx).GetFunction("reducedot")
	}

	reducedot_arg_x1 = x1
	reducedot_arg_x2 = x2
	reducedot_arg_dst = dst
	reducedot_arg_initVal = initVal
	reducedot_arg_n = n

	args := reducedot_argptr[:]
	cu.LaunchKernel(reducedot_code, gridDim.X, gridDim.Y, gridDim.Z, blockDim.X, blockDim.Y, blockDim.Z, 0, reducedot_stream, args)
	reducedot_stream.Synchronize()
	reducedot_lock.Unlock()
}

const reducedot_ptx = `
.version 3.1
.target sm_30
.address_size 64


.visible .entry reducedot(
	.param .u64 reducedot_param_0,
	.param .u64 reducedot_param_1,
	.param .u64 reducedot_param_2,
	.param .f32 reducedot_param_3,
	.param .u32 reducedot_param_4
)
{
	.reg .pred 	%p<8>;
	.reg .s32 	%r<39>;
	.reg .f32 	%f<31>;
	.reg .s64 	%rd<16>;
	// demoted variable
	.shared .align 4 .b8 __cuda_local_var_33846_32_non_const_sdata[2048];

	ld.param.u64 	%rd5, [reducedot_param_0];
	ld.param.u64 	%rd6, [reducedot_param_1];
	ld.param.u64 	%rd7, [reducedot_param_2];
	ld.param.f32 	%f30, [reducedot_param_3];
	ld.param.u32 	%r9, [reducedot_param_4];
	cvta.to.global.u64 	%rd1, %rd7;
	cvta.to.global.u64 	%rd2, %rd6;
	cvta.to.global.u64 	%rd3, %rd5;
	.loc 2 9 1
	mov.u32 	%r38, %ntid.x;
	mov.u32 	%r10, %ctaid.x;
	mov.u32 	%r2, %tid.x;
	mad.lo.s32 	%r37, %r38, %r10, %r2;
	mov.u32 	%r11, %nctaid.x;
	mul.lo.s32 	%r4, %r38, %r11;
	.loc 2 9 1
	setp.ge.s32 	%p1, %r37, %r9;
	@%p1 bra 	BB0_2;

BB0_1:
	.loc 2 9 1
	mul.wide.s32 	%rd8, %r37, 4;
	add.s64 	%rd9, %rd3, %rd8;
	add.s64 	%rd10, %rd2, %rd8;
	ld.global.f32 	%f5, [%rd10];
	ld.global.f32 	%f6, [%rd9];
	fma.rn.f32 	%f30, %f6, %f5, %f30;
	add.s32 	%r37, %r37, %r4;
	.loc 2 9 1
	setp.lt.s32 	%p2, %r37, %r9;
	@%p2 bra 	BB0_1;

BB0_2:
	.loc 2 9 1
	mul.wide.s32 	%rd11, %r2, 4;
	mov.u64 	%rd12, __cuda_local_var_33846_32_non_const_sdata;
	add.s64 	%rd4, %rd12, %rd11;
	st.shared.f32 	[%rd4], %f30;
	bar.sync 	0;
	.loc 2 9 1
	setp.lt.u32 	%p3, %r38, 66;
	@%p3 bra 	BB0_6;

BB0_3:
	.loc 2 9 1
	mov.u32 	%r7, %r38;
	shr.u32 	%r38, %r7, 1;
	.loc 2 9 1
	setp.ge.u32 	%p4, %r2, %r38;
	@%p4 bra 	BB0_5;

	.loc 2 9 1
	ld.shared.f32 	%f7, [%rd4];
	add.s32 	%r16, %r38, %r2;
	mul.wide.u32 	%rd13, %r16, 4;
	add.s64 	%rd15, %rd12, %rd13;
	ld.shared.f32 	%f8, [%rd15];
	add.f32 	%f9, %f7, %f8;
	st.shared.f32 	[%rd4], %f9;

BB0_5:
	.loc 2 9 1
	bar.sync 	0;
	.loc 2 9 1
	setp.gt.u32 	%p5, %r7, 131;
	@%p5 bra 	BB0_3;

BB0_6:
	.loc 2 9 1
	setp.gt.s32 	%p6, %r2, 31;
	@%p6 bra 	BB0_8;

	.loc 2 9 1
	ld.volatile.shared.f32 	%f10, [%rd4];
	ld.volatile.shared.f32 	%f11, [%rd4+128];
	add.f32 	%f12, %f10, %f11;
	st.volatile.shared.f32 	[%rd4], %f12;
	ld.volatile.shared.f32 	%f13, [%rd4+64];
	ld.volatile.shared.f32 	%f14, [%rd4];
	add.f32 	%f15, %f14, %f13;
	st.volatile.shared.f32 	[%rd4], %f15;
	ld.volatile.shared.f32 	%f16, [%rd4+32];
	ld.volatile.shared.f32 	%f17, [%rd4];
	add.f32 	%f18, %f17, %f16;
	st.volatile.shared.f32 	[%rd4], %f18;
	ld.volatile.shared.f32 	%f19, [%rd4+16];
	ld.volatile.shared.f32 	%f20, [%rd4];
	add.f32 	%f21, %f20, %f19;
	st.volatile.shared.f32 	[%rd4], %f21;
	ld.volatile.shared.f32 	%f22, [%rd4+8];
	ld.volatile.shared.f32 	%f23, [%rd4];
	add.f32 	%f24, %f23, %f22;
	st.volatile.shared.f32 	[%rd4], %f24;
	ld.volatile.shared.f32 	%f25, [%rd4+4];
	ld.volatile.shared.f32 	%f26, [%rd4];
	add.f32 	%f27, %f26, %f25;
	st.volatile.shared.f32 	[%rd4], %f27;

BB0_8:
	.loc 2 9 1
	setp.ne.s32 	%p7, %r2, 0;
	@%p7 bra 	BB0_10;

	.loc 2 9 1
	ld.shared.f32 	%f28, [__cuda_local_var_33846_32_non_const_sdata];
	.loc 3 1844 5
	atom.global.add.f32 	%f29, [%rd1], %f28;

BB0_10:
	.loc 2 10 2
	ret;
}


`
