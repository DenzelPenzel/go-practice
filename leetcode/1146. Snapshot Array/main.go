package leetcode

type SnapshotArray struct {
	id   int
	data [][][]int
}

func Constructor(length int) SnapshotArray {
	return SnapshotArray{
		data: make([][][]int, length),
		id:   0,
	}
}

func (sa *SnapshotArray) Set(index int, val int) {
	if len(sa.data[index]) != 0 && sa.data[index][len(sa.data[index])-1][0] == sa.id {
		sa.data[index][len(sa.data[index])-1][1] = val
		return
	}
	sa.data[index] = append(sa.data[index], []int{sa.id, val})
}

func (sa *SnapshotArray) Snap() int {
	sa.id++
	return sa.id - 1
}

func (sa *SnapshotArray) Get(index int, snap_id int) int {
	snap := sa.data[index]

	if len(snap) == 0 {
		return 0
	}

	if snap[0][0] > snap_id {
		return 0
	}

	lo, hi := 0, len(snap)

	for lo < hi {
		mid := lo + ((hi - lo) >> 1)

		if snap[mid][0] <= snap_id {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return snap[lo-1][1]
}

/**
 * Your SnapshotArray object will be instantiated and called as such:
 * obj := Constructor(length);
 * obj.Set(index,val);
 * param_2 := obj.Snap();
 * param_3 := obj.Get(index,snap_id);
 */
