# Code Citations

## License: BSD_2_Clause
https://github.com/sauerbraten/graph/tree/33bdc0ae3dcdf4fdcc490e628daddaab10346fe9/min-priorityqueue.go

```
j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func
```


## License: BSD_2_Clause
https://github.com/andybons/gogif/tree/16d573594812bc09bc62ad1d8a4129c7ba885dc6/mediancut.go

```
j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x
```


## License: BSD_3_Clause
https://github.com/mutecomm/mute/tree/8121965e6779441e0d83339ef41182d792f4d967/msg/session_test.go

```
priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
```


## License: unknown
https://github.com/tricksrobo/naga/tree/33cb03249d8ea69912b20f77da5572520a6ab606/Godeps/_workspace/src/github.com/RoaringBitmap/roaring/priorityqueue.go

```
int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n :
```


## License: unknown
https://github.com/dyrkin/rezerwacje-duw-go/tree/cc373051aaef0598aa9ddc6b95771f8c83630b0b/queue/queue.go

```
) int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i]
```

