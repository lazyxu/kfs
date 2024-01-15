import { View } from 'react-native';
import { Checkbox } from 'react-native-paper';

export default function ({ label, status, onPress }) {
    return (
        <View >
            <Checkbox.Item style={{ padding: 0, margin: 0, paddingLeft: 12, width: "100%" }} label={label} mode='ios'
                status={status}
                onPress={onPress}
            />
        </View>
    )
};
