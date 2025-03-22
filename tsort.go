package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrIncompletePair = errors.New("input contains an odd number of tokens")
	ErrCircular       = errors.New("cycle detected")
)

type successor struct {
	suc  *item
	next *successor
}

type item struct {
	str     string
	left    *item
	right   *item
	balance int8
	printed bool
	count   int
	qlink   *item
	top     *successor
}

var (
	root   = newItem("") // sentinel root for the balanced tree
	zeros  *item
	head   *item
	loop   *item
	nItems int
)

func newItem(str string) *item {
	return &item{str: str}
}

func searchItem(r *item, s string) *item {
	if r.right == nil {
		r.right = newItem(s)
		return r.right
	}
	t := r
	saved := r.right
	p := saved
	sNode := p
	for {
		cmp := strings.Compare(s, p.str)
		if cmp == 0 {
			return p
		}
		var q *item
		if cmp < 0 {
			q = p.left
		} else {
			q = p.right
		}
		if q == nil {
			q = newItem(s)
			if cmp < 0 {
				p.left = q
			} else {
				p.right = q
			}
			cmp2 := strings.Compare(s, sNode.str)
			var rNode *item
			a := 0
			if cmp2 < 0 {
				rNode = sNode.left
				a = -1
			} else {
				rNode = sNode.right
				a = 1
			}
			p2 := rNode
			for p2 != q {
				cmp3 := strings.Compare(s, p2.str)
				if cmp3 < 0 {
					p2.balance = -1
					p2 = p2.left
				} else {
					p2.balance = 1
					p2 = p2.right
				}
			}
			if sNode.balance == 0 || sNode.balance == -int8(a) {
				sNode.balance += int8(a)
				return q
			}
			if rNode.balance == int8(a) {
				p2 = rNode
				if a < 0 {
					sNode.left = rNode.right
					rNode.right = sNode
				} else {
					sNode.right = rNode.left
					rNode.left = sNode
				}
				sNode.balance = 0
				rNode.balance = 0
			} else {
				if a < 0 {
					p2 = rNode.right
					rNode.right = p2.left
					p2.left = rNode
					sNode.left = p2.right
					p2.right = sNode
				} else {
					p2 = rNode.left
					rNode.left = p2.right
					p2.right = rNode
					sNode.right = p2.left
					p2.left = sNode
				}
				sNode.balance = 0
				rNode.balance = 0
				if p2.balance == int8(a) {
					sNode.balance = int8(-a)
				} else if p2.balance == int8(-a) {
					rNode.balance = int8(a)
				}
				p2.balance = 0
			}
			if sNode == t.right {
				t.right = p2
			} else {
				t.left = p2
			}
			return q
		}
		if q.balance != 0 {
			t = p
			sNode = q
		}
		p = q
	}
}

func recordRelation(j, k *item) {
	if j.str != k.str {
		k.count++
		p := &successor{suc: k, next: j.top}
		j.top = p
	}
}

func countItems(k *item) bool {
	nItems++
	return false
}

func scanZeros(k *item) bool {
	if k.count == 0 && !k.printed {
		if head == nil {
			head = k
		} else {
			zeros.qlink = k
		}
		zeros = k
	}
	return false
}

func detectLoop(k *item) bool {
	if k.count > 0 {
		if loop == nil {
			loop = k
		} else {
			p := &k.top
			for *p != nil {
				if (*p).suc == loop {
					if k.qlink != nil {
						tmp := k
						for loop != nil {
							t2 := loop.qlink
							fmt.Fprintln(os.Stderr, loop.str)
							if loop == tmp {
								s := *p
								s.suc.count--
								*p = s.next
								break
							}
							loop.qlink = nil
							loop = t2
						}
						for loop != nil {
							t2 := loop.qlink
							loop.qlink = nil
							loop = t2
						}
						return true
					} else {
						k.qlink = loop
						loop = k
						break
					}
				}
				p = &(*p).next
			}
		}
	}
	return false
}

func recurseTree(r *item, action func(*item) bool) bool {
	if r.left == nil && r.right == nil {
		return action(r)
	}
	if r.left != nil {
		if recurseTree(r.left, action) {
			return true
		}
	}
	if action(r) {
		return true
	}
	if r.right != nil {
		if recurseTree(r.right, action) {
			return true
		}
	}
	return false
}

func walkTree(r *item, action func(*item) bool) {
	if r.right != nil {
		recurseTree(r.right, action)
	}
}

func tsort(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	var j *item
	for scanner.Scan() {
		tok := scanner.Text()
		k := searchItem(root, tok)
		if j != nil {
			recordRelation(j, k)
			k = nil
		}
		j = k
	}

	if j != nil {
		return ErrIncompletePair
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read error: %v", err)
	}

	walkTree(root, countItems)
	ok := true

	for nItems > 0 {
		zeros = nil
		head = nil
		walkTree(root, scanZeros)

		for head != nil {
			p := head.top
			fmt.Println(head.str)
			head.printed = true
			nItems--
			for p != nil {
				p.suc.count--
				if p.suc.count == 0 {
					zeros.qlink = p.suc
					zeros = p.suc
				}
				p = p.next
			}
			head = head.qlink
		}

		if nItems > 0 {
			fmt.Fprintln(os.Stderr, "input contains a loop:")
			ok = false
			for {
				walkTree(root, detectLoop)
				if loop == nil {
					break
				}
			}
		}
	}

	if !ok {
		return ErrCircular
	}

	return nil
}
