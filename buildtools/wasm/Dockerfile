FROM rust:1-slim

RUN apt-get update && apt-get install -y binaryen curl
RUN rustup target add wasm32-wasip1-threads
ADD buildtools/wasm/version.txt version.txt

RUN mkdir -p /workspace && curl -L https://github.com/owenlamont/ryl/archive/refs/tags/$(cat version.txt | awk '{$1=$1};1').tar.gz | tar -xz --strip-components 1 -C /workspace
WORKDIR /workspace

RUN cargo update && cargo build --release --target wasm32-wasip1-threads --no-default-features

SHELL ["/bin/bash", "-c"]
RUN for f in target/wasm32-wasip1-threads/release/*.wasm; do wasm-opt -o "target/$(basename $f)" --flatten --rereloop --converge -O3 "$f"; done

CMD ["bash", "-c", "cp target/*.wasm /out/"]
