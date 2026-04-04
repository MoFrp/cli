#!/bin/bash
set -e

# 接受版本号参数，如果没有则尝试从 git tag 或 version 文件获取
VERSION=${1:-}
if [ -z "$VERSION" ]; then
    VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
fi
if [ -z "$VERSION" ]; then
    VERSION=$(cat ./version 2>/dev/null || echo "0.68.1")
fi

echo "Package version: $VERSION"

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd openbsd android'
arch_all='amd64 arm arm64 mips64 mips64le mips mipsle riscv64 loong64'
extra_all='_ hf'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        for extra in $extra_all; do
            suffix="${os}_${arch}"
            if [ "x${extra}" != x"_" ]; then
                suffix="${os}_${arch}_${extra}"
            fi
            frp_dir_name="frpc_${VERSION}_${suffix}"
            frp_path="./packages/frpc_${VERSION}_${suffix}"

            if [ "x${os}" = x"windows" ]; then
                if [ ! -f "./frpc_${os}_${arch}.exe" ]; then
                    continue
                fi
                mkdir -p ${frp_path}
                mv ./frpc_${os}_${arch}.exe ${frp_path}/frpc.exe
            else
                if [ ! -f "./frpc_${suffix}" ]; then
                    continue
                fi
                mkdir -p ${frp_path}
                mv ./frpc_${suffix} ${frp_path}/frpc
            fi
            cp ../LICENSE ${frp_path}
            cp -f ../conf/frpc.toml ${frp_path}

            # packages
            cd ./packages
            if [ "x${os}" = x"windows" ]; then
                zip -rq ${frp_dir_name}.zip ${frp_dir_name}
            else
                tar -zcf ${frp_dir_name}.tar.gz ${frp_dir_name}
            fi
            cd ..
            rm -rf ${frp_path}
        done
    done
done

cd -
echo "Package done!"
