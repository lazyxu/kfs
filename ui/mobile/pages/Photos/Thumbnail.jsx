import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { Pressable } from 'react-native';
import FastImage from './FastImage';

export default function ({ hash, width, navigation }) {
    // const { ref, inView } = useInView({ threshold: 0 });
    return (
        <Pressable onPress={() => navigation.navigate("Viewer", { hash })}  >
            <FastImage style={{
                    width: width,
                    height: width,
                }} source={{ uri: `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}` }} 
            />
        </Pressable>
    );
}
