import color from 'color';
import { useState } from 'react';
import { Pressable, View } from 'react-native';
import { Icon, Text, useTheme } from 'react-native-paper';

{/* <Button contentStyle={contentStyle} labelStyle={labelStyle} icon="file-image-outline" onPress={() => navigation.navigate("SearchType")}></Button> */ }
const iconSize = 18;
export default function ({ children, icon, style, labelStyle, onPress, disabled }) {
    const [isActive, setIsActive] = useState(false);
    const { colors } = useTheme();
    // console.log("theme", theme)
    const activeStyle = isActive && { backgroundColor: color(colors.onSurface).alpha(0.12).rgb().string() };
    const textStyle = {
        fontWeight: 500,
        lineHeight: 20,
        fontSize: 14,
        marginLeft: 16,
        marginRight: 16,
        marginTop: 9,
        marginBottom: 9,
        color: disabled ? colors.onSurfaceDisabled : colors.primary,
        flex: 1,
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-between",
    }
    return (
        <View style={{ width: "100%", ...activeStyle }}>
            <Pressable style={{ width: "100%" }} onPressIn={() => setIsActive(true)} onPressOut={() => setIsActive(false)} onPress={onPress} disabled={disabled}>
                <View style={{
                    display: "flex",
                    flexDirection: "row",
                    alignItems: "center",
                    borderRadius: 20,
                    width: "100%",
                    ...style
                }}>
                    <View style={{
                        marginLeft: 12,
                        marginRight: -4,
                    }}>
                        <Icon
                            source={icon}
                            size={iconSize}
                            color={disabled ? colors.onSurfaceDisabled : colors.primary}
                        />
                    </View>
                    {typeof children === 'string' ?
                        <Text style={textStyle}>
                            {children}
                        </Text> :
                        <View style={textStyle}>
                            {children}
                        </View>
                    }
                </View>
            </Pressable>
        </View>
    )
};
