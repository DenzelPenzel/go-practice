```

func newproc(size int32, fn *funcval) {
 argp := add(unsafe.Pointer(&fn), sys.PtrSize)
 gp := getg()
 pc := getcallerpc()

 systemstack(func() {
  newg := newproc1(fn, argp, size, gp, pc)

  _p_ := getg().m.p.ptr()

  // newly created Goroutine will call this method to decide how to schedule
  runqput(_p_, newg, true)

  if mainStarted {
   wakep()
  }
 })

}
 ...

 if next {
  retryNext:
   oldnext := _p_.runnext
   // When next is true, the new Goroutine will always be put into the next scheduling field
   if !_p_.runnext.cas(oldnext, guintptr(unsafe.Pointer(gp))) {
    goto retryNext
   }

   if oldnext == 0 {
    return
   }

   // Kick the old runnext out to the regular run queue
   gp := oldnext.ptr()

 }

```
