import React from 'react';
import { View, Text, StyleSheet, FlatList } from 'react-native';

const Leaderboard = ({ users }) => {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Лидеры</Text>
      <FlatList
        data={users}
        keyExtractor={(item) => item.id.toString()}
        renderItem={({ item }) => (
          <View style={styles.userItem}>
            <Text style={styles.position}>{item.position}</Text>
            <View style={styles.userInfo}>
              <Text style={styles.name}>{item.name}</Text>
              <Text style={styles.points}>{item.points} баллов</Text>
            </View>
          </View>
        )}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    backgroundColor: '#f5f5f5',
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 20,
    color: '#2c3e50',
  },
  userItem: {
    flexDirection: 'row',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderColor: '#eee',
  },
  position: {
    width: 40,
    textAlign: 'center',
    color: '#2ecc71',
  },
  userInfo: {
    flex: 1,
    marginLeft: 10,
  },
  name: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#2c3e50',
  },
  points: {
    color: '#666',
  },
});

export default Leaderboard;