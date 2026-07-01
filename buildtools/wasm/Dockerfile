FROM rust:1-slim AS wasm

RUN apt-get update && apt-get install -y binaryen curl wabt
ADD buildtools/wasm/version.txt version.txt

RUN mkdir -p /workspace && curl -L https://github.com/owenlamont/ryl/archive/refs/tags/$(cat version.txt | awk '{$1=$1};1').tar.gz | tar -xz --strip-components 1 -C /workspace
WORKDIR /workspace

RUN rustup target add wasm32-wasip1-threads

ENV RUSTFLAGS="-C link-arg=--export=__wasilibc_cwd -C link-arg=--export=malloc"
RUN cargo update && cargo build --release --target wasm32-wasip1-threads --no-default-features

SHELL ["/bin/bash", "-c"]
RUN for f in target/wasm32-wasip1-threads/release/*.wasm; do wasm-opt -o "target/$(basename $f)" --flatten --rereloop --converge -O3 "$f"; done

# Emit <name>.memory.go holding initialPages, read from each module's declared
# memory minimum; the host allocates that many pages up front (see host.go).
RUN for f in target/*.wasm; do \
      pages=$(wasm-objdump -x "$f" | grep -oE 'initial=[0-9]+' | head -1 | grep -oE '[0-9]+'); \
      printf '// Code generated from %s. DO NOT EDIT.\n\npackage wasm2go\n\nconst initialPages = %s\n' "$(basename "$f")" "$pages" > "${f%.wasm}.memory.go"; \
    done

FROM golang:1.26-trixie

COPY --from=wasm /workspace/target/*.wasm /build/
COPY --from=wasm /workspace/target/*.memory.go /build/
SHELL ["/bin/bash", "-c"]
RUN for f in /build/*.wasm; do go run github.com/ncruces/wasm2go@v0.4.10 -unsafe < "$f" > "$f".go; sed -i '1i //go:build !codeanalysis' "$f".go; done

CMD ["bash", "-c", "cp /build/*.go /out/"]
