# Path to this plugin
PROTOC_GEN_TS_PATH="./ui/node_modules/.bin/protoc-gen-ts"

OUTPUT_DIR=./ui/src/pb
protoc -I=pb fs.proto \
    --plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
    --js_out="import_style=commonjs,binary:${OUTPUT_DIR}" \
    --ts_out=service=grpc-web:"${OUTPUT_DIR}"
sed -i "" "1i\\"$'\n'" /* eslint-disable */"$'\n' ${OUTPUT_DIR}/fs_pb.js
