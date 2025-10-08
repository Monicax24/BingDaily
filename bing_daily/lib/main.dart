import 'package:bing_daily/constants.dart';
import 'package:bing_daily/firebase_options.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:bing_daily/utils/route_utils.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(
    name: 'bing-daily',
    options: DefaultFirebaseOptions.currentPlatform,
  );
  runApp(const ProviderScope(child: MyApp()));
}

/// Root widget for the Bing Daily app with routing and theming.
class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return MaterialApp.router(
      title: appTitle,
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: bingGreen,
          primary: bingGreen,
          secondary: bingAccent,
        ),
        scaffoldBackgroundColor: bingWhite,
      ),
      routerConfig: createRouter(ref),
    );
  }
}
