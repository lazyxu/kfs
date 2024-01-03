import { View } from 'react-native';
import { IconButton, Surface } from "react-native-paper";

export default function ({ navigation, hash }) {
    return (
        <Surface style={{
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: "space-between",
            paddingHorizontal: 4,
            position: "absolute", left: 0, top: 0, right: 0,
            zIndex: 1,
        }}>
            <IconButton
                // size={size}
                onPress={() => navigation.pop()}
                // iconColor={actionIconColor}
                icon="chevron-left"
            // disabled={disabled}
            // rippleColor={rippleColor}
            />
            <View style={{
                flexDirection: 'row',
                alignItems: 'center',
            }}>
                <IconButton
                    // size={size}
                    // onPress={downloadImage}
                    // iconColor={actionIconColor}
                    icon="download-outline"
                // disabled={disabled}
                // rippleColor={rippleColor}
                />
                <IconButton
                    // size={size}
                    // onPress={downloadImage}
                    // iconColor={actionIconColor}
                    icon="dots-vertical"
                // disabled={disabled}
                // rippleColor={rippleColor}
                />
            </View>
        </Surface>
    );
}
