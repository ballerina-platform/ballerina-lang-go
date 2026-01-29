package semtypes

import (
	"ballerina-lang-go/common"
	"iter"
)

type subtypePairIterator struct {
	i1              int
	i2              int
	t1              []BasicSubtype
	t2              []BasicSubtype
	bits            BasicTypeBitSet
	doneIteration   bool
	shouldCalculate bool
	cache           *SubtypePair
}

func (i *subtypePairIterator) include(code BasicTypeCode) bool {
	return (i.bits.bitset & (1 << code.Code)) != 0
}

func (i *subtypePairIterator) get1() BasicSubtype {
	return i.t1[i.i1]
}

func (i *subtypePairIterator) get2() BasicSubtype {
	return i.t2[i.i2]
}

func (i *subtypePairIterator) hasNext() bool {
	if i.doneIteration {
		return false
	}
	if i.shouldCalculate {
		cache := i.internalNext()
		if cache == nil {
			i.doneIteration = true
		}
		i.cache = cache
		i.shouldCalculate = false
	}
	return !i.doneIteration
}

func (i *subtypePairIterator) next() SubtypePair {
	if i.doneIteration {
		panic("Exhausted iterator")
	}
	if i.shouldCalculate {
		cache := i.internalNext()
		if cache == nil {
			panic("unexpected nil cache")
		}
		i.cache = cache
	}
	i.shouldCalculate = true
	return *i.cache
}

func (i *subtypePairIterator) internalNext() *SubtypePair {
	for {
		if i.i1 >= len(i.t1) {
			if i.i2 >= len(i.t2) {
				break
			}
			t := i.get2()
			code := t.BasicTypeCode
			data2 := t.SubtypeData
			i.i2++
			if i.include(code) {
				return common.ToPointer(CreateSubTypePair(code, nil, data2))
			}
		} else if i.i2 >= len(i.t2) {
			t := i.get1()
			code := t.BasicTypeCode
			data1 := t.SubtypeData
			i.i1++
			if i.include(code) {
				return common.ToPointer(CreateSubTypePair(code, data1, nil))
			}
		} else {
			t1 := i.get1()
			code1 := t1.BasicTypeCode
			data1 := t1.SubtypeData

			t2 := i.get2()
			code2 := t2.BasicTypeCode
			data2 := t2.SubtypeData
			if code1 == code2 {
				i.i1++
				i.i2++
				if i.include(code1) {
					return common.ToPointer(CreateSubTypePair(code1, data1, data2))
				}
			} else if code1.Code < code2.Code {
				i.i1++
				if i.include(code1) {
					return common.ToPointer(CreateSubTypePair(code1, data1, nil))
				}
			} else {
				i.i2++
				if i.include(code2) {
					return common.ToPointer(CreateSubTypePair(code2, nil, data2))
				}
			}
		}
	}
	return nil
}

func (i *subtypePairIterator) toIterator() iter.Seq[SubtypePair] {
	return func(yield func(SubtypePair) bool) {
		for i.hasNext() {
			if !yield(i.next()) {
				break
			}
		}
	}
}

func newSubtypePairs(s1, s2 SemType, bits BasicTypeBitSet) iter.Seq[SubtypePair] {
	i := &subtypePairIterator{
		i1:              0,
		i2:              0,
		t1:              unpackToBasicSubtypes(s1),
		t2:              unpackToBasicSubtypes(s2),
		bits:            bits,
		doneIteration:   false,
		shouldCalculate: true,
		cache:           nil,
	}
	return i.toIterator()
}

func unpackToBasicSubtypes(t SemType) []BasicSubtype {
	if _, ok := t.(*BasicTypeBitSet); ok {
		return nil
	}
	return unpackComplexSemType(t.(ComplexSemType))
}
