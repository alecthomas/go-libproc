package libproc

// #include <libproc.h>
// #include <golibproc.h>
// #include <unistd.h>
import "C"
import (
	"fmt"
	"unsafe"
)

var RawHighestRUsageFlavor = RawDiscoverRUsageVersion()

// RawProcListPidsPath is a low-level binding for proc_listpidspath.
func RawProcListPidsPath(_type uint32, typeInfo uint32, path string, pathFlags uint32, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.listpidspath(C.uint32_t(_type), C.uint32_t(typeInfo), C.CString(path), C.uint32_t(pathFlags), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcListPids is a low-level binding for proc_listpids.
func RawProcListPids(_type uint32, typeInfo uint32, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.listpids(C.uint32_t(_type), C.uint32_t(typeInfo), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcListAllPids is a low-level binding for proc_listallpids.
func RawProcListAllPids(buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.listallpids(buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcListPGRPids is a low-level binding for proc_listpgrppids.
func RawProcListPGRPids(pgrpid Pid, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.listpgrppids(C.pid_t(pgrpid), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcListChildPids is a low-level binding for proc_listchildpids.
func RawProcListChildPids(ppid Pid, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.listchildpids(C.pid_t(ppid), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcPidInfo is a low-level binding for proc_pidinfo.
func RawProcPidInfo(pid Pid, flavor PidInfoFlavor, arg uint64, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.pidinfo(C.int(pid), C.int(flavor), C.uint64_t(arg), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcPidFDInfo is a low-level binding for proc_pidfdinfo.
func RawProcPidFDInfo(pid Pid, fd int, flavor PidFdInfoFlavor, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.pidfdinfo(C.int(pid), C.int(fd), C.int(flavor), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcPidFileReportInfo is a low-level binding for proc_pidfileportinfo.
func RawProcPidFileReportInfo(pid Pid, fileport uint32, flavor PidFileReportInfoFlavor, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.pidfileportinfo(C.int(fileport), C.uint32_t(fileport), C.int(flavor), buffer, C.int(bufferSize))
	return int(n), getErr(err)
}

// RawProcName is a low-level binding for proc_name.
func RawProcName(pid Pid, buffer unsafe.Pointer, bufferSize uint32) (int, error) {
	n, err := C.name(C.int(pid), buffer, C.uint32_t(bufferSize))
	return int(n), getErr(err)
}

// RawRegionFilename is a low-level binding for proc_regionfilename.
func RawRegionFilename(pid Pid, address uint64, buffer unsafe.Pointer, bufferSize int) (int, error) {
	n, err := C.regionfilename(C.int(pid), C.uint64_t(address), buffer, C.uint32_t(bufferSize))
	return int(n), getErr(err)
}

// RawKMsgBuf is a low-level binding for proc_kmsgbuf.
func RawKMsgBuf(buffer unsafe.Pointer, bufferSize uint32) (int, error) {
	n, err := C.kmsgbuf(buffer, C.uint32_t(bufferSize))
	return int(n), getErr(err)
}

// RawPidPath is a low-level binding for proc_pidpath.
func RawPidPath(pid Pid, buffer unsafe.Pointer, bufferSize uint32) (int, error) {
	n, err := C.pidpath(C.int(pid), buffer, C.uint32_t(bufferSize))
	return int(n), getErr(err)
}

// RawLibVersion is a low-level binding for proc_libversion.
func RawLibVersion(major *int, minor *int) (int, error) {
	n, err := C.libversion((*C.int)(unsafe.Pointer(major)), (*C.int)(unsafe.Pointer(minor)))
	return int(n), getErr(err)
}

// RawDiscoverRUsageVersion finds the highest supported version of the rusage_info struct at runtime.
//
// Will return -1 on error.
func RawDiscoverRUsageVersion() RUsageFlavor {
	pid := C.getpid()
	var v4 C.struct_rusage_info_v4
	for i := 4; i >= 0; i-- {
		_, err := C.pid_rusage(C.int(pid), C.int(i), (*C.rusage_info_t)(unsafe.Pointer(&v4)))
		if err == nil {
			return RUsageFlavor(i)
		}
		fmt.Println(err)
	}
	return -1
}

// RawPidRUsage is a low-level binding for proc_pid_rusage.
func RawPidRUsage(pid Pid) (*RUsage, int, error) {
	var rusage_info C.struct_rusage_info_v4
	n, err := C.pid_rusage(C.int(pid), C.int(RawHighestRUsageFlavor), (*C.rusage_info_t)(unsafe.Pointer(&rusage_info)))
	if err != nil {
		return nil, int(n), getErr(err)
	}
	rusage := &RUsage{
		UserTime:                  uint64(rusage_info.ri_user_time),
		SystemTime:                uint64(rusage_info.ri_system_time),
		PkgIdleWkups:              uint64(rusage_info.ri_pkg_idle_wkups),
		InterruptWkups:            uint64(rusage_info.ri_interrupt_wkups),
		Pageins:                   uint64(rusage_info.ri_pageins),
		WiredSize:                 uint64(rusage_info.ri_wired_size),
		ResidentSize:              uint64(rusage_info.ri_resident_size),
		PhysFootprint:             uint64(rusage_info.ri_phys_footprint),
		ProcStartAbstime:          uint64(rusage_info.ri_proc_start_abstime),
		ProcExitAbstime:           uint64(rusage_info.ri_proc_exit_abstime),
		ChildUserTime:             uint64(rusage_info.ri_child_user_time),
		ChildSystemTime:           uint64(rusage_info.ri_child_system_time),
		ChildPkgIdleWkups:         uint64(rusage_info.ri_child_pkg_idle_wkups),
		ChildInterruptWkups:       uint64(rusage_info.ri_child_interrupt_wkups),
		ChildPageins:              uint64(rusage_info.ri_child_pageins),
		ChildElapsedAbstime:       uint64(rusage_info.ri_child_elapsed_abstime),
		DiskIoBytesRead:           uint64(rusage_info.ri_diskio_bytesread),
		DiskIoBytesWritten:        uint64(rusage_info.ri_diskio_byteswritten),
		CPUTimeQosDefault:         uint64(rusage_info.ri_cpu_time_qos_default),
		CPUTimeQosMaintenance:     uint64(rusage_info.ri_cpu_time_qos_maintenance),
		CPUTimeQosBackground:      uint64(rusage_info.ri_cpu_time_qos_background),
		CPUTimeQosUtility:         uint64(rusage_info.ri_cpu_time_qos_utility),
		CPUTimeQosLegacy:          uint64(rusage_info.ri_cpu_time_qos_legacy),
		CPUTimeQosUserInitiated:   uint64(rusage_info.ri_cpu_time_qos_user_initiated),
		CPUTimeQosUserInteractive: uint64(rusage_info.ri_cpu_time_qos_user_interactive),
		BilledSystemTime:          uint64(rusage_info.ri_billed_system_time),
		ServicedSystemTime:        uint64(rusage_info.ri_serviced_system_time),
		LogicalWrites:             uint64(rusage_info.ri_logical_writes),
		LifetimeMaxPhysFootprint:  uint64(rusage_info.ri_lifetime_max_phys_footprint),
		Instructions:              uint64(rusage_info.ri_instructions),
		Cycles:                    uint64(rusage_info.ri_cycles),
		BilledEnergy:              uint64(rusage_info.ri_billed_energy),
		ServicedEnergy:            uint64(rusage_info.ri_serviced_energy),
	}
	for i := 0; i < 16; i++ {
		rusage.UUID[i] = byte(rusage_info.ri_uuid[i])
	}
	return rusage, 0, nil
}

// RawSetPControl is a low-level binding for proc_setpcontrol.
func RawSetPControl(control ProcessControl) (int, error) {
	n, err := C.setpcontrol(C.int(control))
	return int(n), getErr(err)
}

// RawProcTrackDirty is a low-level binding for proc_track_dirty.
func RawProcTrackDirty(pid Pid, flags uint32) (int, error) {
	n, err := C.track_dirty(C.pid_t(pid), C.uint32_t(flags))
	return int(n), getErr(err)
}

// RawProcSetDirty is a low-level binding for proc_set_dirty.
func RawProcSetDirty(pid Pid, dirty bool) (int, error) {
	n, err := C.set_dirty(C.pid_t(pid), C._Bool(dirty))
	return int(n), getErr(err)
}

// RawProcGetDirty is a low-level binding for proc_get_dirty.
func RawProcGetDirty(pid Pid, flags *uint32) (int, error) {
	n, err := C.get_dirty(C.pid_t(pid), (*C.uint32_t)(unsafe.Pointer(flags)))
	return int(n), getErr(err)
}

// RawProcClearDirty is a low-level binding for proc_clear_dirty.
func RawProcClearDirty(pid Pid, flags uint32) (int, error) {
	n, err := C.clear_dirty(C.pid_t(pid), C.uint32_t(flags))
	return int(n), getErr(err)
}
