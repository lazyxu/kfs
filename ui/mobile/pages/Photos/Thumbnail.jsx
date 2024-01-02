import ImgCancelable from '@kfs/common/components/ImgCancelable';
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useState } from "react";
import { Image, Pressable, View } from 'react-native';

export default function ({ hash, width, navigation, index, list }) {
    const [] = useState(false);
    const [inView, setInView] = useState(false);
    const src = `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`;
    // console.log('Thumbnail.width:', width);
    // console.log('Thumbnail.Inview:', inView);
    useEffect(() => {
        setInView(true);
    }, []);
    return (
        <Pressable onPress={() => navigation.navigate("Viewer", { hash, index, list })}>
            {/* <InView onChange={(inView) => { setInView(inView); }}> */}
                <ImgCancelable inView={inView}
                    src={src}
                    renderImg={(url) => <Image style={{
                        width: width,
                        height: width,
                    }} source={{ uri: url }} />}
                    renderSkeleton={() => <View style={{
                        width: width,
                        height: width,
                    }} />}
                />
            {/* </InView> */}
        </Pressable>
    );
}
