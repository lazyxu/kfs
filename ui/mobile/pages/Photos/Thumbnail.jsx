import ImgCancelable from '@kfs/common/components/ImgCancelable';
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useInView } from "react-intersection-observer";
import { Image, Pressable, View } from 'react-native';

export default function ({ hash, width, navigation }) {
    const { ref, inView } = useInView({ threshold: 0 });
    return (
        <Pressable ref={ref} onPress={() => navigation.navigate("Viewer", { hash })}  >
            <ImgCancelable inView={inView}
                src={`${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`}
                renderImg={(url) => <Image style={{
                    width: width,
                    height: width,
                }} source={{ uri: url }} />}
                renderSkeleton={() => <View style={{
                    width: width,
                    height: width,
                }} />}
            />
        </Pressable>
    );
}
