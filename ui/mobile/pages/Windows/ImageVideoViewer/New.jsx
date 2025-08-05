import ImageView from "@staltz/react-native-image-viewing";
import { useState } from "react";
import { Text } from "react-native-paper";
import { View } from "react-native-web";
import { shareAsync } from 'expo-sharing';
import { Platform } from 'react-native';
import { IconButton, Surface } from "react-native-paper";
import DCIMAttributeRN from "@kfs/common/components/File/DCIMAttributeRN";

export default function ({ navigation, route }) {
  const { list, index } = route.params;
  const [currentImageIndex, setCurrentImageIndex] = useState(index);
  const [showDCIMAttribute, setShowDCIMAttribute] = useState(false);
  return (
    <>
      {showDCIMAttribute && <DCIMAttributeRN
        dcimFile={list[currentImageIndex]}
        onClose={() => setShowDCIMAttribute(false)}
      />}
      {!showDCIMAttribute && <ImageView
        images={list}
        imageIndex={index}
        visible={true}
        onRequestClose={() => navigation.pop()}
        onImageIndexChange={setCurrentImageIndex}
        // 自定义工具栏
        FooterComponent={({ imageIndex }) => (
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
                  onPress={async () => {
                    try {
                      let uri = list[imageIndex].uri;
                      console.log('shareImage', Platform.OS, uri);
                      await shareAsync(uri);
                    } catch (e) {
                      console.log(e);
                    }
                  }}
                  icon="export-variant"
                  disabled={!list[imageIndex].uri}
              />
              <IconButton
                  style={{ borderRadius: 0, backgroundColor: null, transition: null }}
                  onPress={() => setShowDCIMAttribute(true)}
                  icon="information-outline"
              />
              <IconButton
                  icon="trash-can-outline"
                  disabled={true}
              />
              <View style={{ padding: 20 }}>
                <Text>图片 {imageIndex + 1}/{list.length}</Text>
              </View>
          </Surface>
        )}
      />}
    </>
  );
}
