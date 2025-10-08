import 'package:flutter/material.dart';

class CameraPage extends StatelessWidget {
  const CameraPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Bing Daily')),
      body: const Center(child: Text('Welcome to the camera Page!')),
    );
  }
}
