# 版权 @2023 nes-wa 作者。保留所有权利。

run:
	wa run -target=js -console

.PHONY: neswa
neswa:
	-rm -rf ./src/nes
	make -C ../neswago/nes
	mv ../neswago/nes/output ./src/nes
	-rm ./src/nes/zz_helper.wa ./src/nes/static.wa 
	cp zz_helper.wa.txt ./src/nes/zz_helper.wa

.PHONY: nesgo
nesgo:
	make -C ./nesgo

clean:
	-rm -rf ./output
