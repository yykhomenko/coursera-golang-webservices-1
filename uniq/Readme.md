Scans input and prints unique lines only. Input data must be sorted, otherwise panic occurs. 

Usage:
```bash
cat data.txt | go run main.go
```
or:
```bash
echo -e "1\n2\n3\n4\n4\n5" | go run main.go
```
```bash
mac:w1-uniq yykhomenko$ echo -e "1\n2\n3\n4\n4\n5" | go run main.go 
1
2
3
4
5

```
