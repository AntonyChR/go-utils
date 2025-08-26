package assert

import "testing"

func TestAssert(t *testing.T) {
  Assert(true,"the tested value is false", t)
}

func TestAssertEq(t *testing.T) {
  AssertEq(true,true,"the tested values ar not equal", t)
}

func TestAssertNe(t *testing.T) {
  AssertNe(true,false,"the tested values are equal", t)
}
