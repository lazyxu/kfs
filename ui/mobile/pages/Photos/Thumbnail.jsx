import ImgCancelable from '@kfs/common/components/ImgCancelable';
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useInView } from "react-intersection-observer";
import { Pressable, View } from 'react-native';
import FastImage from './FastImage';

export default function ({ hash, width, navigation, index, list }) {
    const { ref, inView } = useInView({ threshold: 0 });
    const src = `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`;
    return (
        <Pressable ref={ref} onPress={() => navigation.navigate("Viewer", { hash, index, list })}  >
            <ImgCancelable inView={inView}
                src={src}
                renderImg={(url) => <FastImage style={{
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
