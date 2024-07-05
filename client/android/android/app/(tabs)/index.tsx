// app/(tabs)/index.tsx
import React from 'react';
import { SafeAreaView, StyleSheet, Text } from 'react-native';
import Autocomplete from '../Autocomplete';

const App = () => {
  return (
    <SafeAreaView style={styles.container}>
      <Text style={styles.title}>Restaurant Autocomplete</Text>
      <Autocomplete />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 40,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 20,
  },
});

export default App;
