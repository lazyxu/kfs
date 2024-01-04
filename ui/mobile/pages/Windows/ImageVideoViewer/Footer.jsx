import { shareAsync } from 'expo-sharing';
import { Platform } from 'react-native';
import { IconButton, Surface } from "react-native-paper";

export default function ({ navigation, hash, uri }) {
    const shareImage = async () => {
        try {
            console.log('shareImage', Platform.OS, uri);
            await shareAsync(uri);
        } catch (e) {
            console.log(e);
        }
    };
    return (
        <Surface style={{
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: "space-around",
            paddingHorizontal: 4,
            position: "absolute", left: 0, bottom: 0, right: 0,
            zIndex: 1,
        }}>
            <IconButton
                style={{ borderRadius: 0, backgroundColor: null, transition: null }}
                onPress={shareImage}
                icon="export-variant"
                disabled={!uri}
            />
            <IconButton
                style={{ borderRadius: 0, backgroundColor: null, transition: null }}
                onPress={() => navigation.navigate("Info", { hash })}
                icon="information-outline"
            />
            <IconButton
                icon="trash-can-outline"
                disabled={true}
            />
        </Surface>
    );
}
