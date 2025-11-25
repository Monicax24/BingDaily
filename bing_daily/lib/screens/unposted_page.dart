// import 'package:bing_daily/constants.dart';
// import 'package:flutter/material.dart';

// /// Home page displaying a welcome message for the Bing Daily app.
// class UnpostedPage extends StatelessWidget {
//   const UnpostedPage({super.key});

//   @override
//   Widget build(BuildContext context) {
//     return Scaffold(
//       appBar: AppBar(
//         title: const Text(
//           "Prompt of the Day:",
//           style: TextStyle(color: Colors.white),
//         ),
//         backgroundColor: bingGreen,
//       ),
//       body: Column(
//         mainAxisAlignment: MainAxisAlignment.end,
//         crossAxisAlignment: CrossAxisAlignment.center,
//         children: [
//           const Padding(
//             padding: EdgeInsets.all(16.0),
//             child: Text(
//               'Hold up ...',
//               style: TextStyle(
//                 fontSize: 24,
//                 fontWeight: FontWeight.bold,
//                 color: Colors.black87,
//               ),
//               textAlign: TextAlign.center,
//             ),
//           ),
//           SizedBox(height: 26),
//           const Padding(
//             padding: EdgeInsets.all(16.0),
//             child: Text(
//               'Post before you scroll',
//               style: TextStyle(
//                 fontSize: 18,
//                 fontWeight: FontWeight.bold,
//                 color: Colors.black87,
//               ),
//               textAlign: TextAlign.center,
//             ),
//           ),
//           SizedBox(height: 34),
//           Container(
//             height: 460,
//             decoration: BoxDecoration(
//               image: DecorationImage(
//                 image: AssetImage('assets/images/baxter_stopsign.png'),
//                 fit: BoxFit.fitHeight,
//               ),
//             ),
//           ),
//         ],
//       ),
//       backgroundColor: bingWhite,
//     );
//   }
// }
