import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';

/// Home page displaying a welcome message for the Bing Daily app.
class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text(appTitle), backgroundColor: bingGreen),
      body: const Center(
        child: Text(
          'Welcome to the Home Page!',
          style: TextStyle(color: bingGreen, fontSize: 20),
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
