import useSysConfig from '@kfs/common/hox/sysConfig';
import { View } from 'react-native';
import { Appbar, Button, Divider, RadioButton, Text, TextInput } from 'react-native-paper';

export default () => {
    const { sysConfig, setSysConfig, resetSysConfig } = useSysConfig();
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="更多" />
            </Appbar.Header>
            <View>
                <Text>设置：</Text>
                <Button mode="elevated" style={{ width: 200 }} onPress={e => resetSysConfig()}>恢复默认设置</Button>
                <Divider />
                <View>
                    <Text>主题：</Text>
                    <RadioButton.Group
                        value={sysConfig.theme}
                        onValueChange={theme => setSysConfig(c => ({ ...c, theme }))}
                    >
                        {["light", "dark", "system"].map(value =>
                            <View key={value}>
                                <RadioButton.Item key={value} value={value} label={value} />
                            </View>
                        )}
                    </RadioButton.Group>
                </View>
                <Divider />
                <View>
                    <TextInput mode="outlined" label="Web服务器"
                        value={sysConfig.webServer}
                        onChangeText={webServer => setSysConfig(c => ({ ...c, webServer }))}
                    />
                </View>
                <Divider />
                <View>
                    <TextInput mode="outlined" label="Socket服务器"
                        value={sysConfig.socketServer}
                        onChangeText={socketServer => setSysConfig(c => ({ ...c, socketServer }))}
                    />
                </View>
                <Divider />
                <View>
                    <TextInput mode="outlined" label="客户端Web服务器端口"
                        value={sysConfig.port}
                        onChangeText={port => setSysConfig(c => ({ ...c, port }))}
                    />
                </View>
            </View>
        </>
    );
};
