package flow

import (
	"github.com/chrislusf/gleam/instruction"
)

// GroupBy e.g. GroupBy(Field(1,2,3)) group data by field 1,2,3
func (d *Dataset) GroupBy(name string, sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	ret := d.LocalSort(name, sortOption).LocalGroupBy(name, sortOption)
	if len(d.Shards) > 1 {
		ret = ret.MergeSortedTo(name, 1, sortOption).LocalGroupBy(name, sortOption)
	}
	ret.IsLocalSorted = sortOption.orderByList
	return ret
}

func (d *Dataset) LocalGroupBy(name string, sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	ret, step := add1ShardTo1Step(d)
	indexes := sortOption.Indexes()
	ret.IsPartitionedBy = indexes
	step.SetInstruction(name, instruction.NewLocalGroupBySorted(sortOption.Indexes()))
	return ret
}
