000000000046fb80 <main.main>:
  46fb80:	4c 8d 64 24 c8       	lea    r12,[rsp-0x38]
  46fb85:	4d 3b 66 10          	cmp    r12,QWORD PTR [r14+0x10]
  46fb89:	0f 86 b0 01 00 00    	jbe    46fd3f <main.main+0x1bf>
  46fb8f:	55                   	push   rbp
  46fb90:	48 89 e5             	mov    rbp,rsp
  46fb93:	48 81 ec b0 00 00 00 	sub    rsp,0xb0
  46fb9a:	48 8d 05 7f bf 00 00 	lea    rax,[rip+0xbf7f]        # 47bb20 <type:*+0xbb20>
  46fba1:	31 db                	xor    ebx,ebx
  46fba3:	b9 a0 86 01 00       	mov    ecx,0x186a0
  46fba8:	e8 73 7e ff ff       	call   467a20 <runtime.makeslice>
  46fbad:	31 c9                	xor    ecx,ecx
  46fbaf:	31 d2                	xor    edx,edx
  46fbb1:	31 db                	xor    ebx,ebx
  46fbb3:	be a0 86 01 00       	mov    esi,0x186a0
  46fbb8:	eb 17                	jmp    46fbd1 <main.main+0x51>
  46fbba:	4d 8d 48 ff          	lea    r9,[r8-0x1]
  46fbbe:	49 c1 e1 04          	shl    r9,0x4
  46fbc2:	4a 89 0c 08          	mov    QWORD PTR [rax+r9*1],rcx
  46fbc6:	4a 89 4c 08 08       	mov    QWORD PTR [rax+r9*1+0x8],rcx
  46fbcb:	48 ff c1             	inc    rcx
  46fbce:	4c 89 c3             	mov    rbx,r8
  46fbd1:	48 81 f9 a0 86 01 00 	cmp    rcx,0x186a0
  46fbd8:	73 75                	jae    46fc4f <main.main+0xcf>
  46fbda:	4c 8d 43 01          	lea    r8,[rbx+0x1]
  46fbde:	66 90                	xchg   ax,ax
  46fbe0:	4c 39 c6             	cmp    rsi,r8
  46fbe3:	73 d5                	jae    46fbba <main.main+0x3a>
  46fbe5:	48 89 4c 24 48       	mov    QWORD PTR [rsp+0x48],rcx
  46fbea:	88 54 24 47          	mov    BYTE PTR [rsp+0x47],dl
  46fbee:	49 83 f8 02          	cmp    r8,0x2
  46fbf2:	7f 2f                	jg     46fc23 <main.main+0xa3>
  46fbf4:	84 d2                	test   dl,dl
  46fbf6:	75 2b                	jne    46fc23 <main.main+0xa3>
  46fbf8:	48 85 db             	test   rbx,rbx
  46fbfb:	75 26                	jne    46fc23 <main.main+0xa3>
  46fbfd:	44 0f 11 bc 24 88 00 	movups XMMWORD PTR [rsp+0x88],xmm15
  46fc04:	00 00
  46fc06:	44 0f 11 bc 24 98 00 	movups XMMWORD PTR [rsp+0x98],xmm15
  46fc0d:	00 00
  46fc0f:	48 8d 84 24 88 00 00 	lea    rax,[rsp+0x88]
  46fc16:	00
  46fc17:	be 02 00 00 00       	mov    esi,0x2
  46fc1c:	ba 01 00 00 00       	mov    edx,0x1
  46fc21:	eb 97                	jmp    46fbba <main.main+0x3a>
  46fc23:	4c 89 c3             	mov    rbx,r8
  46fc26:	48 89 f1             	mov    rcx,rsi
  46fc29:	bf 01 00 00 00       	mov    edi,0x1
  46fc2e:	48 8d 35 eb be 00 00 	lea    rsi,[rip+0xbeeb]        # 47bb20 <type:*+0xbb20>
  46fc35:	e8 c6 7e ff ff       	call   467b00 <runtime.growslice>
  46fc3a:	49 89 d8             	mov    r8,rbx
  46fc3d:	48 89 ce             	mov    rsi,rcx
  46fc40:	0f b6 54 24 47       	movzx  edx,BYTE PTR [rsp+0x47]
  46fc45:	48 8b 4c 24 48       	mov    rcx,QWORD PTR [rsp+0x48]
  46fc4a:	e9 6b ff ff ff       	jmp    46fbba <main.main+0x3a>
  46fc4f:	48 89 5c 24 50       	mov    QWORD PTR [rsp+0x50],rbx
  46fc54:	48 89 84 24 a8 00 00 	mov    QWORD PTR [rsp+0xa8],rax
  46fc5b:	00
  46fc5c:	48 8d 05 fd 60 00 00 	lea    rax,[rip+0x60fd]        # 475d60 <type:*+0x5d60>
  46fc63:	bb a0 86 01 00       	mov    ebx,0x186a0
  46fc68:	48 89 d9             	mov    rcx,rbx
  46fc6b:	e8 b0 7d ff ff       	call   467a20 <runtime.makeslice>
  46fc70:	48 8b 54 24 50       	mov    rdx,QWORD PTR [rsp+0x50]
  46fc75:	48 8b b4 24 a8 00 00 	mov    rsi,QWORD PTR [rsp+0xa8]
  46fc7c:	00
  46fc7d:	31 c9                	xor    ecx,ecx
  46fc7f:	31 db                	xor    ebx,ebx
  46fc81:	bf a0 86 01 00       	mov    edi,0x186a0
  46fc86:	41 b8 a0 86 01 00    	mov    r8d,0x186a0
  46fc8c:	eb 12                	jmp    46fca0 <main.main+0x120>
  46fc8e:	4e 89 5c c8 f8       	mov    QWORD PTR [rax+r9*8-0x8],r11
  46fc93:	49 8d 4a 01          	lea    rcx,[r10+0x1]
  46fc97:	4c 89 cf             	mov    rdi,r9
  46fc9a:	66 0f 1f 44 00 00    	nop    WORD PTR [rax+rax*1+0x0]
  46fca0:	48 39 d1             	cmp    rcx,rdx
  46fca3:	0f 8d 8d 00 00 00    	jge    46fd36 <main.main+0x1b6>
  46fca9:	4c 8d 4f 01          	lea    r9,[rdi+0x1]
  46fcad:	49 89 ca             	mov    r10,rcx
  46fcb0:	48 c1 e1 04          	shl    rcx,0x4
  46fcb4:	4c 8b 1c 0e          	mov    r11,QWORD PTR [rsi+rcx*1]
  46fcb8:	4d 39 c8             	cmp    r8,r9
  46fcbb:	73 d1                	jae    46fc8e <main.main+0x10e>
  46fcbd:	4c 89 54 24 60       	mov    QWORD PTR [rsp+0x60],r10
  46fcc2:	4c 89 5c 24 58       	mov    QWORD PTR [rsp+0x58],r11
  46fcc7:	88 5c 24 46          	mov    BYTE PTR [rsp+0x46],bl
  46fccb:	49 83 f9 04          	cmp    r9,0x4
  46fccf:	7f 27                	jg     46fcf8 <main.main+0x178>
  46fcd1:	84 db                	test   bl,bl
  46fcd3:	75 23                	jne    46fcf8 <main.main+0x178>
  46fcd5:	48 85 ff             	test   rdi,rdi
  46fcd8:	75 1e                	jne    46fcf8 <main.main+0x178>
  46fcda:	44 0f 11 7c 24 68    	movups XMMWORD PTR [rsp+0x68],xmm15
  46fce0:	44 0f 11 7c 24 78    	movups XMMWORD PTR [rsp+0x78],xmm15
  46fce6:	48 8d 44 24 68       	lea    rax,[rsp+0x68]
  46fceb:	41 b8 04 00 00 00    	mov    r8d,0x4
  46fcf1:	bb 01 00 00 00       	mov    ebx,0x1
  46fcf6:	eb 96                	jmp    46fc8e <main.main+0x10e>
  46fcf8:	4c 89 cb             	mov    rbx,r9
  46fcfb:	4c 89 c1             	mov    rcx,r8
  46fcfe:	bf 01 00 00 00       	mov    edi,0x1
  46fd03:	48 8d 35 56 60 00 00 	lea    rsi,[rip+0x6056]        # 475d60 <type:*+0x5d60>
  46fd0a:	e8 f1 7d ff ff       	call   467b00 <runtime.growslice>
  46fd0f:	48 8b 54 24 50       	mov    rdx,QWORD PTR [rsp+0x50]
  46fd14:	48 8b b4 24 a8 00 00 	mov    rsi,QWORD PTR [rsp+0xa8]
  46fd1b:	00
  46fd1c:	4c 8b 54 24 60       	mov    r10,QWORD PTR [rsp+0x60]
  46fd21:	4c 8b 5c 24 58       	mov    r11,QWORD PTR [rsp+0x58]
  46fd26:	49 89 d9             	mov    r9,rbx
  46fd29:	49 89 c8             	mov    r8,rcx
  46fd2c:	0f b6 5c 24 46       	movzx  ebx,BYTE PTR [rsp+0x46]
  46fd31:	e9 58 ff ff ff       	jmp    46fc8e <main.main+0x10e>
  46fd36:	48 81 c4 b0 00 00 00 	add    rsp,0xb0
  46fd3d:	5d                   	pop    rbp
  46fd3e:	c3                   	ret
  46fd3f:	90                   	nop
  46fd40:	e8 3b ad ff ff       	call   46aa80 <runtime.morestack_noctxt.abi0>
  46fd45:	e9 36 fe ff ff       	jmp    46fb80 <main.main>
  46fd4a:	cc                   	int3
  46fd4b:	cc                   	int3
  46fd4c:	cc                   	int3
  46fd4d:	cc                   	int3
  46fd4e:	cc                   	int3
  46fd4f:	cc                   	int3
  46fd50:	cc                   	int3
  46fd51:	cc                   	int3
  46fd52:	cc                   	int3
  46fd53:	cc                   	int3
  46fd54:	cc                   	int3
  46fd55:	cc                   	int3
  46fd56:	cc                   	int3
  46fd57:	cc                   	int3
  46fd58:	cc                   	int3
  46fd59:	cc                   	int3
  46fd5a:	cc                   	int3
  46fd5b:	cc                   	int3
  46fd5c:	cc                   	int3
  46fd5d:	cc                   	int3
  46fd5e:	cc                   	int3
  46fd5f:	cc                   	int3
