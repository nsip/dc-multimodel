
cd tools; go build release.go; cd ..
./tools/release dc-multimodel > version/version.go
sh build.sh
./tools/release dc-multimodel dc-multimodel-Mac.zip build/dc-multimodel-Mac.zip
./tools/release dc-multimodel dc-multimodel-Win64.zip build/dc-multimodel-Win64.zip
./tools/release dc-multimodel dc-multimodel-Win32.zip build/dc-multimodel-Win32.zip
./tools/release dc-multimodel dc-multimodel-Linux64.zip build/dc-multimodel-Linux64.zip
./tools/release dc-multimodel dc-multimodel-Linux32.zip build/dc-multimodel-Linux32.zip
