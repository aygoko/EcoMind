import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Define the Bloc
abstract class AppEvent {}

class AppBloc extends Bloc<AppEvent, String> {
  AppBloc() : super("Open up App.js to start working on your app!") {
    // No event handling needed for this basic example
  }
}

// App entry point
class MyApp extends StatelessWidget {
  @override
 Wich build(BuildContext context) {
    return BlocProvider(
      create: (context) => AppBloc(),
      child: MaterialApp(
        home: HomeScreen(),
      ),
    );
  }
}

// Main screen widget
class HomeScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: BlocBuilder<AppBloc, String>(
          builder: (context, state) {
            return Text(state);
          },
        ),
      ),
    );
  }
}

void main() {
  runApp(MyApp());
}