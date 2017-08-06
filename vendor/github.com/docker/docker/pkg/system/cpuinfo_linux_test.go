package system

import (
	"strings"
	"testing"
)

func TestCpuInfo(t *testing.T) {
	const input = `processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 0
cpu cores	: 6
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 1
cpu cores	: 6
apicid		: 2
initial apicid	: 2
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 2
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 2
cpu cores	: 6
apicid		: 4
initial apicid	: 4
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 3
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 3
cpu cores	: 6
apicid		: 6
initial apicid	: 6
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 4
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 4
cpu cores	: 6
apicid		: 8
initial apicid	: 8
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 5
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 5
cpu cores	: 6
apicid		: 10
initial apicid	: 10
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 6
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.730
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 0
cpu cores	: 6
apicid		: 32
initial apicid	: 32
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 7
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 1
cpu cores	: 6
apicid		: 34
initial apicid	: 34
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 8
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 2
cpu cores	: 6
apicid		: 36
initial apicid	: 36
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 9
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 3
cpu cores	: 6
apicid		: 38
initial apicid	: 38
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 10
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 4
cpu cores	: 6
apicid		: 40
initial apicid	: 40
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 11
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 5
cpu cores	: 6
apicid		: 42
initial apicid	: 42
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 12
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 0
cpu cores	: 6
apicid		: 1
initial apicid	: 1
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 13
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 1
cpu cores	: 6
apicid		: 3
initial apicid	: 3
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 14
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 2
cpu cores	: 6
apicid		: 5
initial apicid	: 5
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 15
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 3
cpu cores	: 6
apicid		: 7
initial apicid	: 7
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 16
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 4
cpu cores	: 6
apicid		: 9
initial apicid	: 9
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 17
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 0
siblings	: 12
core id		: 5
cpu cores	: 6
apicid		: 11
initial apicid	: 11
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4599.91
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 18
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 0
cpu cores	: 6
apicid		: 33
initial apicid	: 33
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 19
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.359
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 1
cpu cores	: 6
apicid		: 35
initial apicid	: 35
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 20
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.730
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 2
cpu cores	: 6
apicid		: 37
initial apicid	: 37
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 21
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 3
cpu cores	: 6
apicid		: 39
initial apicid	: 39
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 22
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2300.000
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 4
cpu cores	: 6
apicid		: 41
initial apicid	: 41
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

processor	: 23
vendor_id	: GenuineIntel
cpu family	: 6
model		: 45
model name	: Intel(R) Xeon(R) CPU E5-2630 0 @ 2.30GHz
stepping	: 7
microcode	: 0x710
cpu MHz		: 2299.910
cache size	: 15360 KB
physical id	: 1
siblings	: 12
core id		: 5
cpu cores	: 6
apicid		: 43
initial apicid	: 43
fpu		: yes
fpu_exception	: yes
cpuid level	: 13
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 ds_cpl vmx smx est tm2 ssse3 cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic popcnt tsc_deadline_timer aes xsave avx lahf_lm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid xsaveopt
bogomips	: 4604.78
clflush size	: 64
cache_alignment	: 64
address sizes	: 46 bits physical, 48 bits virtual
power management:

`
	cpuinfo, err := parseCpuInfo(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if cpuinfo.Processors != 24 {
		t.Fatalf("Unexpected total processors, expected: 24, got: %v", cpuinfo.Processors)
	}

	if len(cpuinfo.Sockets) != 2 {
		t.Fatalf("Unexpected total sockets, expected: 2, got: %v", len(cpuinfo.Sockets))
	}

	for _, v := range [][4]int64{
		{0, 0, 0, 12},
		{0, 1, 1, 13},
		{0, 2, 2, 14},
		{0, 3, 3, 15},
		{0, 4, 4, 16},
		{0, 5, 5, 17},
		{1, 0, 6, 18},
		{1, 1, 7, 19},
		{1, 2, 8, 20},
		{1, 3, 9, 21},
		{1, 4, 10, 22},
		{1, 5, 11, 23},
	} {
		socket := cpuinfo.Sockets.Get(v[0])
		if socket == nil {
			t.Fatalf("Unexpected socket, cannot find socket %v", v[0])
		}

		core := socket.Cores.Get(v[1])
		if socket == nil {
			t.Fatalf("Unexpected core, cannot find core %v in socket %v", v[1], v[0])
		}

		if core.Processors[0] != v[2] || core.Processors[1] != v[3] {
			t.Fatalf(
				"Unexpected processor of core %v in socket %v, expected: [%v, %v], got: %v",
				v[1], v[0], v[2], v[3], core.Processors,
			)
		}
	}
}
