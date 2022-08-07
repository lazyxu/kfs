# npm install -g protoc-gen-ts google-protobuf

OUTPUT_DIR=./ui/src
protoc pb/fs.proto \
    --js_out="import_style=commonjs,binary:${OUTPUT_DIR}" \
    --ts_out=service=grpc-web:"${OUTPUT_DIR}"
# sed -i "" "1i\\"$'\n'" /* eslint-disable */"$'\n' ${OUTPUT_DIR}/fs_pb.js
# sed -i "" "1i\\"$'\n'" /* eslint-disable */"$'\n' ${OUTPUT_DIR}/fs_pb_service.js
