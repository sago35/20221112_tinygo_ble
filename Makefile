smoketest: FORCE
	go build -o ${TMP}/tmp ./central1
	go build -o ${TMP}/tmp ./central2
	tinygo build -target xiao-ble -o ${TMP}/tmp.uf2 ./central1
	tinygo build -target xiao-ble -o ${TMP}/tmp.uf2 ./central2
	tinygo build -target xiao-ble -o ${TMP}/tmp.uf2 ./peripheral1
	tinygo build -target xiao-ble -o ${TMP}/tmp.uf2 ./peripheral2
	tinygo build -target xiao-ble -o ${TMP}/tmp.uf2 ./peripheral2b

FORCE:
