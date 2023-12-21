import useSysConfig from '@kfs/common/hox/sysConfig';
import { View } from 'react-native';
import { Appbar, Button, RadioButton, Text } from 'react-native-paper';

export default () => {
    const { sysConfig, setSysConfig, resetSysConfig } = useSysConfig();
    return (
        <>
            <Appbar.Header>
                <Appbar.Content title="更多" />
            </Appbar.Header>
            <View>
                <Text>设置：</Text>
                <Button mode="elevated" style={{ width: 200 }} onPress={e => resetSysConfig()}>恢复默认设置</Button>
                <View>
                    <Text>主题：</Text>
                    <RadioButton.Group
                        value={sysConfig.theme}
                        onValueChange={theme => setSysConfig(c => ({ ...c, theme }))}
                    >
                        {["light", "dark", "system"].map(value =>
                            <View>
                                <RadioButton.Item key={value} value={value} label={value} />
                            </View>
                        )}
                    </RadioButton.Group>
                </View>
            </View>
        </>
    );
};
