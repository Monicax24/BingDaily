import 'dart:io' show Platform;
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';

class NotificationApi {
  final _firebaseMessaging = FirebaseMessaging.instance;

  /// Initializes notifications by requesting permission and fetching the FCM token with retry on iOS.
  Future<void> initNotifications() async {
    await _firebaseMessaging.requestPermission();
    final token = await _getFCMTokenWithRetry();
    if (token != null) {
      print('Token: $token');
    } else {
      print('Failed to retrieve FCM token after retries');
    }
  }

  /// Retrieves FCM token with retry mechanism for iOS to handle APNS token delay.
  Future<String?> _getFCMTokenWithRetry({
    int retries = 3,
    Duration delay = const Duration(seconds: 2),
  }) async {
    for (int attempt = 0; attempt < retries; attempt++) {
      try {
        return await _firebaseMessaging.getToken();
      } catch (e) {
        if (e is FirebaseException &&
            e.code == 'apns-token-not-set' &&
            Platform.isIOS) {
          if (attempt < retries - 1) {
            await Future.delayed(delay);
            continue;
          }
        }
        print('Error getting FCM token: $e');
        return null;
      }
    }
    return null;
  }
}
