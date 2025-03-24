import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import CircularProgress from '../components/CircularProgress';
import Leaderboard from '../components/Leaderboard'; 

const DashboardScreen = () => {
  const progress = 60;
  const goal = 100;

  const users = [
    { id: 1, name: 'Иван', points: 300, position: 1 },
    { id: 2, name: 'Мария', points: 250, position: 2 },
    { id: 3, name: 'Алексей', points: 200, position: 3 },
  ];

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Главная страница</Text>
      <View style={styles.chartContainer}>
        <CircularProgress progress={progress} goal={goal} />
        <Text style={styles.progressText}>
          {progress}% из {goal}%
        </Text>
      </View>
      <Leaderboard users={users} /> 
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
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 20,
    color: '#2c3e50',
  },
  chartContainer: {
    alignItems: 'center',
    marginBottom: 20,
  },
  progressText: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2ecc71',
  },
});

export default DashboardScreen;