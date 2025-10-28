import 'package:bing_daily/models/app_user.dart';
import 'package:bing_daily/providers/auth_provider.dart';
import 'package:bing_daily/screens/login_screen.dart';
import 'package:bing_daily/screens/username_setup_screen.dart';
import 'package:bing_daily/widgets/main_screen.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

/// Defines the Go Router configuration for the Bing Daily app.
GoRouter createRouter(WidgetRef ref) {
  return GoRouter(
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const LoginScreen(),
        redirect: (context, state) {
          final AppUser? user = ref.watch(userNotifierProvider);
          if (user != null) {
            return '/main'; // Redirect to main if logged in
          }
          return null;
        },
      ),
      GoRoute(
        path: '/main',
        builder: (context, state) => const MainScreen(),
        redirect: (context, state) {
          final AppUser? user = ref.watch(userNotifierProvider);
          if (user == null) {
            return '/'; // Redirect to login if not logged in
          }
          return null;
        },
      ),
      GoRoute(
        path: '/username',
        builder: (context, state) => const UsernameSetupScreen(),
      ),
    ],
  );
}
