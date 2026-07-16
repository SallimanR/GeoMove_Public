gps_update=api/protobuf/geolocation/gps_realtime_channel/gps_update.proto
gps_get=api/protobuf/geolocation/gps_realtime_channel/gps_get.proto
ws_message=api/protobuf/websocket/ws_message.proto

protoc --go_out=backend/monolith/internal/domains/geolocation/interface/websocket/pb $gps_update
protoc --go_out=backend/monolith/internal/domains/geolocation/interface/websocket/pb $gps_get
protoc --go_out=backend/monolith/internal/websockethub/proto -I api/protobuf/websocket $ws_message

pnpm --dir=frontend exec protoc \
	--ts_proto_out=packages/@domains/geolocation/src/api/proto \
	--ts_proto_opt=useOptionals=none,esModuleInterop=true \
	-I ../api/protobuf/geolocation/gps_realtime_channel \
	../api/protobuf/geolocation/gps_realtime_channel/gps_update.proto

pnpm --dir=frontend exec protoc \
	--ts_proto_out=packages/@domains/geolocation/src/api/proto \
	--ts_proto_opt=useOptionals=none,esModuleInterop=true \
	-I ../api/protobuf/geolocation/gps_realtime_channel \
	../api/protobuf/geolocation/gps_realtime_channel/gps_get.proto

pnpm --dir=frontend exec protoc \
	--ts_proto_out=packages/@domains/geolocation/src/api/proto \
	--ts_proto_opt=useOptionals=none,esModuleInterop=true \
	-I ../api/protobuf/websocket \
	../api/protobuf/websocket/ws_message.proto
