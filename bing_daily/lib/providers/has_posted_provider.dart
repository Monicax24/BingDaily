import 'package:flutter_riverpod/legacy.dart';

/// Provider to simulate if the user has posted today (true/false).
final hasPostedProvider = StateProvider<bool>((ref) => false);
