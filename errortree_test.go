package errortree_test

import (
	"errors"
	"testing"

	"github.com/jackc/errortree"
	"github.com/stretchr/testify/require"
)

func TestNodeAddGetSelf(t *testing.T) {
	tree := &errortree.Node{}
	tree.Add(nil, errors.New("foo"))
	errs := tree.Get(nil)
	require.Len(t, errs, 1)
	require.EqualError(t, errs[0], "foo")

	tree.Add(nil, errors.New("bar"))
	errs = tree.Get(nil)
	require.Len(t, errs, 2)
	require.EqualError(t, errs[0], "foo")
	require.EqualError(t, errs[1], "bar")
}

func TestNodeAddGetAttributes(t *testing.T) {
	tree := &errortree.Node{}
	tree.Add([]any{"name"}, errors.New("foo"))
	errs := tree.Get([]any{"name"})
	require.Len(t, errs, 1)
	require.EqualError(t, errs[0], "foo")

	tree.Add([]any{"name"}, errors.New("bar"))
	errs = tree.Get([]any{"name"})
	require.Len(t, errs, 2)
	require.EqualError(t, errs[0], "foo")
	require.EqualError(t, errs[1], "bar")

	tree.Add([]any{"age"}, errors.New("quz"))
	errs = tree.Get([]any{"age"})
	require.Len(t, errs, 1)
	require.EqualError(t, errs[0], "quz")

	errs = tree.Get([]any{"nonexistent", "branch"})
	require.Len(t, errs, 0)

	tree.Add([]any{"abc", "def", "ghi"}, errors.New("deep error"))
	errs = tree.Get([]any{"abc", "def", "ghi"})
	require.Len(t, errs, 1)
	require.EqualError(t, errs[0], "deep error")
}

func TestNodeAddGetElements(t *testing.T) {
	tree := &errortree.Node{}
	tree.Add([]any{13}, errors.New("foo"))
	errs := tree.Get([]any{13})
	require.Len(t, errs, 1)
	require.EqualError(t, errs[0], "foo")

	tree.Add([]any{13}, errors.New("bar"))
	errs = tree.Get([]any{13})
	require.Len(t, errs, 2)
	require.EqualError(t, errs[0], "foo")
	require.EqualError(t, errs[1], "bar")

	errs = tree.Get([]any{7})
	require.Len(t, errs, 0)
}

func TestNodeAddNodeMergesErrors(t *testing.T) {
	tree := &errortree.Node{}
	tree.Add([]any{"abc"}, errors.New("foo"))

	node := &errortree.Node{}
	node.Add([]any{}, errors.New("bar"))
	node.Add([]any{}, errors.New("baz"))

	tree.Add([]any{"abc"}, node)

	errs := tree.Get([]any{"abc"})
	require.Len(t, errs, 3)
	require.EqualError(t, errs[0], "foo")
	require.EqualError(t, errs[1], "bar")
	require.EqualError(t, errs[2], "baz")
}
