import { createStackNavigator } from '@react-navigation/stack';
import { NavigationContainer } from '@react-navigation/native';

import RegistrationScreen from '../screens/AuthScreens/RegistrationScreen';
import LoginScreen from '../screens/AuthScreens/LoginScreen';
import HomeScreen from '../screens/HomeScreen';
import TasksScreen from '../screens/TasksScreen';
import StatsScreen from '../screens/StatsScreen';
import ProfileScreen from '../screens/ProfileScreen';
import SurveyScreen from '../screens/SurveyScreen';
import SurveyResultScreen from '../screens/SurveyResultScreen';
import DashboardScreen from '../screens/DashboardScreen';
import WeeklyGoalsScreen from '../screens/WeeklyGoalsScreen';

const Stack = createStackNavigator();

export function AppNavigator() {
  return (
    <NavigationContainer>
      <Stack.Navigator 
        initialRouteName="Login"
        screenOptions={{
          headerShown: false
        }}
      >
        <Stack.Screen name="Registration" component={RegistrationScreen} />
        <Stack.Screen name="Login" component={LoginScreen} />
        <Stack.Screen name="Home" component={HomeScreen} />
        <Stack.Screen name="Tasks" component={TasksScreen} />
        <Stack.Screen name="Stats" component={StatsScreen} />
        <Stack.Screen name="Profile" component={ProfileScreen} />
        <Stack.Screen name="Survey" component={SurveyScreen} />
        <Stack.Screen name="SurveyResult" component={SurveyResultScreen} />
        <Stack.Screen name="Dashboard" component={DashboardScreen} />
        <Stack.Screen name="WeeklyGoals" component={WeeklyGoalsScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}