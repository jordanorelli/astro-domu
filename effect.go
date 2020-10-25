package main

type effect interface {
	apply(*room)
}
