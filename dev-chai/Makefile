# 版权 @2023 nes-wa 作者。保留所有权利。

run:
	wa run

.PHONY: neswa
neswa:
	-rm -rf ././src/nes
	make -C ./neswago/nes
	mv ./neswago/nes/output ./src/nes
	-rm ./src/nes/zz_helper.wa
	cp zz_helper.wa.txt ./src/nes/zz_helper.wa
	#wa run -target=wasi

.PHONY: nesgo
nesgo:
	make -C ./nesgo

wasi:
	-@make clean
	wa run -target=wasi

clean:
	-rm -rf ./output
