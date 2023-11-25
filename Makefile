# 版权 @2023 nes-wa 作者。保留所有权利。

run:
	-@make clean
	wa run

.PHONY: nesgo
nesgo:
	make -C ./nesgo

wasi:
	-@make clean
	wa run -target=wasi

clean:
	-rm -rf ./output
