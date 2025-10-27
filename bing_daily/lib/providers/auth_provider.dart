import 'package:bing_daily/models/app_user.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Provider for Firebase authentication state, tracking user login status.
final authStateProvider = StreamProvider<User?>((ref) {
  return FirebaseAuth.instance.authStateChanges();
});

class userNotifier extends Notifier<AppUser?> {
  @override
  AppUser? build() => null;

  void setUser(AppUser? user) {
    state = user;
  }
}

final userNotifierProvider =
    NotifierProvider<userNotifier, AppUser?>(userNotifier.new);
