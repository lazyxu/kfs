import { IconButton, Surface } from "react-native-paper";

export default function ({ navigation, hash }) {
    return (
        <Surface style={{
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: "space-around",
            paddingHorizontal: 4,
            position: "fixed", left: 0, bottom: 0, right: 0,
            zIndex: 1,
        }}>
            <IconButton
                // size={size}
                // onPress={shareImage}
                // iconColor={actionIconColor}
                icon="export-variant"
            // disabled={disabled}
            // rippleColor={rippleColor}
            />
            <IconButton
                // size={size}
                onPress={() => navigation.navigate("Info", { hash })}
                // iconColor={actionIconColor}
                icon="information-outline"
            // disabled={disabled}
            // rippleColor={rippleColor}
            />
            <IconButton
                // size={size}
                // onPress={downloadImage}
                // iconColor={actionIconColor}
                icon="trash-can-outline"
            // disabled={disabled}
            // rippleColor={rippleColor}
            />
        </Surface>
    );
}
